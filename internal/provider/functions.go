package provider

import (
	"context"
	"fmt"

	"github.com/germanbrew/terraform-provider-dotenv/internal/dotenv"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ function.Function = getFunction{}
	_ function.Function = getByKeyFunction{}
	_ function.Function = hasFunction{}
	_ function.Function = getOrDefaultFunction{}
	_ function.Function = decodeFunction{}
	_ function.Function = encodeFunction{}
)

const (
	descKeyParameter      = "Name of the key"
	descFilenameParameter = "`Default: .env` Path to the dotenv file. " + descRelativePaths
)

func keyParameter() function.StringParameter {
	return function.StringParameter{
		Name:                "key",
		MarkdownDescription: descKeyParameter,
	}
}

func optionalFilenameParameter() function.StringParameter {
	return function.StringParameter{
		Name:                "filename",
		MarkdownDescription: descFilenameParameter + " This argument is optional.",
	}
}

// optionalFilename returns the filename passed to a variadic filename
// parameter, or the default. argPos is the position of the variadic parameter.
func optionalFilename(filenames []string, argPos int64) (string, *function.FuncError) {
	switch len(filenames) {
	case 0:
		return defaultFilename, nil
	case 1:
		if filenames[0] == "" {
			return defaultFilename, nil
		}

		return filenames[0], nil
	default:
		return "", function.NewArgumentFuncError(argPos, "at most one filename can be given")
	}
}

// lookupKey reads filename and returns the value of key and whether it exists.
func lookupKey(ctx context.Context, key, filename string) (string, bool, *function.FuncError) {
	entries, err := readFileEntries(ctx, filename)
	if err != nil {
		return "", false, function.NewFuncError(err.Error())
	}

	value, exists := entries[key]

	return value, exists, nil
}

func getValue(ctx context.Context, key, filename string) (string, *function.FuncError) {
	value, exists, funcErr := lookupKey(ctx, key, filename)
	if funcErr != nil {
		return "", funcErr
	}

	if !exists {
		return "", function.NewFuncError(fmt.Sprintf("Could not find key '%s' in file '%s'", key, filename))
	}

	return value, nil
}

// --- get ---

// NewGetFunction returns the `get` function.
func NewGetFunction() function.Function {
	return getFunction{}
}

type getFunction struct{}

func (f getFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "get"
}

func (f getFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Get the value of a key from a dotenv file",
		MarkdownDescription: "Reads a dotenv file and returns the value of a single entry by its key. " +
			"Fails if the key does not exist; see `get_or_default` and `has` for alternatives.\n\n" +
			descSupportedFormats + "\n\n" + descFunctionExpansion,
		Parameters:        []function.Parameter{keyParameter()},
		VariadicParameter: optionalFilenameParameter(),
		Return:            function.StringReturn{},
	}
}

func (f getFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var (
		key       string
		filenames []string
	)

	resp.Error = req.Arguments.Get(ctx, &key, &filenames)
	if resp.Error != nil {
		return
	}

	filename, funcErr := optionalFilename(filenames, 1)
	if funcErr != nil {
		resp.Error = funcErr

		return
	}

	value, funcErr := getValue(ctx, key, filename)
	if funcErr != nil {
		resp.Error = funcErr

		return
	}

	resp.Error = resp.Result.Set(ctx, value)
}

// --- get_by_key (deprecated) ---

// NewGetByKeyFunction returns the deprecated `get_by_key` function.
func NewGetByKeyFunction() function.Function {
	return getByKeyFunction{}
}

type getByKeyFunction struct{}

func (f getByKeyFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "get_by_key"
}

func (f getByKeyFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Get by key function (deprecated, use get)",
		MarkdownDescription: "~> **Deprecated:** Use [`get`](" + docsBaseURL + "/functions/get) instead. " +
			"`get_by_key(key, filename)` is equivalent to `get(key, filename)` and will be removed in the next major version.\n\n" +
			"Reads and provides a single entry of a .env file by its key\n\n" +
			descSupportedFormats,
		DeprecationMessage: "get_by_key is deprecated and will be removed in the next major version, use get instead.",
		Parameters: []function.Parameter{
			keyParameter(),
			function.StringParameter{
				Name:                "filename",
				MarkdownDescription: "Path to the dotenv file",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f getByKeyFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var key, filename string

	resp.Error = req.Arguments.Get(ctx, &key, &filename)
	if resp.Error != nil {
		return
	}

	tflog.Warn(ctx, "provider::dotenv::get_by_key is deprecated and will be removed in the next major version, use provider::dotenv::get instead")

	filename, _ = optionalFilename([]string{filename}, 1)

	value, funcErr := getValue(ctx, key, filename)
	if funcErr != nil {
		resp.Error = funcErr

		return
	}

	resp.Error = resp.Result.Set(ctx, value)
}

// --- has ---

// NewHasFunction returns the `has` function.
func NewHasFunction() function.Function {
	return hasFunction{}
}

type hasFunction struct{}

func (f hasFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "has"
}

func (f hasFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Check whether a dotenv file contains a key",
		MarkdownDescription: "Reads a dotenv file and returns `true` if it contains the given key.\n\n" + descSupportedFormats,
		Parameters:          []function.Parameter{keyParameter()},
		VariadicParameter:   optionalFilenameParameter(),
		Return:              function.BoolReturn{},
	}
}

func (f hasFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var (
		key       string
		filenames []string
	)

	resp.Error = req.Arguments.Get(ctx, &key, &filenames)
	if resp.Error != nil {
		return
	}

	filename, funcErr := optionalFilename(filenames, 1)
	if funcErr != nil {
		resp.Error = funcErr

		return
	}

	_, exists, funcErr := lookupKey(ctx, key, filename)
	if funcErr != nil {
		resp.Error = funcErr

		return
	}

	resp.Error = resp.Result.Set(ctx, exists)
}

// --- get_or_default ---

// NewGetOrDefaultFunction returns the `get_or_default` function.
func NewGetOrDefaultFunction() function.Function {
	return getOrDefaultFunction{}
}

type getOrDefaultFunction struct{}

func (f getOrDefaultFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "get_or_default"
}

func (f getOrDefaultFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Get the value of a key from a dotenv file, or a default",
		MarkdownDescription: "Reads a dotenv file and returns the value of a single entry by its key, " +
			"or `default` if the key does not exist. The file itself must exist.\n\n" +
			descSupportedFormats + "\n\n" + descFunctionExpansion,
		Parameters: []function.Parameter{
			keyParameter(),
			function.StringParameter{
				Name:                "default",
				MarkdownDescription: "Value returned if the key does not exist",
			},
		},
		VariadicParameter: optionalFilenameParameter(),
		Return:            function.StringReturn{},
	}
}

func (f getOrDefaultFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var (
		key          string
		defaultValue string
		filenames    []string
	)

	resp.Error = req.Arguments.Get(ctx, &key, &defaultValue, &filenames)
	if resp.Error != nil {
		return
	}

	filename, funcErr := optionalFilename(filenames, 2)
	if funcErr != nil {
		resp.Error = funcErr

		return
	}

	value, exists, funcErr := lookupKey(ctx, key, filename)
	if funcErr != nil {
		resp.Error = funcErr

		return
	}

	if !exists {
		value = defaultValue
	}

	resp.Error = resp.Result.Set(ctx, value)
}

// --- decode ---

// NewDecodeFunction returns the `decode` function.
func NewDecodeFunction() function.Function {
	return decodeFunction{}
}

type decodeFunction struct{}

func (f decodeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "decode"
}

func (f decodeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Parse a string in the dotenv format",
		MarkdownDescription: "Parses a string in the dotenv format and returns its entries as a map. " +
			"Useful together with [`file()`](https://developer.hashicorp.com/terraform/language/functions/file), " +
			"[`templatefile()`](https://developer.hashicorp.com/terraform/language/functions/templatefile) " +
			"or values from other providers.\n\n" + descSupportedFormats + "\n\n" +
			"Variable references (`$VAR`, `${VAR}`) that are not defined in the content are resolved from " +
			"the environment of the Terraform process.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "content",
				MarkdownDescription: "Content in the dotenv format",
			},
		},
		Return: function.MapReturn{ElementType: types.StringType},
	}
}

func (f decodeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var content string

	resp.Error = req.Arguments.Get(ctx, &content)
	if resp.Error != nil {
		return
	}

	entries, err := dotenv.Parse(content)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())

		return
	}

	resp.Error = resp.Result.Set(ctx, entries)
}

// --- encode ---

// NewEncodeFunction returns the `encode` function.
func NewEncodeFunction() function.Function {
	return encodeFunction{}
}

type encodeFunction struct{}

func (f encodeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "encode"
}

func (f encodeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Render a map as a string in the dotenv format",
		MarkdownDescription: "Renders a map of strings in the dotenv format, one `KEY=value` line per entry, sorted by key. " +
			"Values are written unquoted when possible, otherwise single or double quoted, so that `decode` and the " +
			"data sources read back exactly the same value (variable references like `$VAR` are not expanded).\n\n" +
			"Keys may only contain letters, digits, underscores and dots. Values that contain a single quote or a " +
			"line break *and* a literal `\\n` or `\\r` sequence can not be represented in the dotenv format and cause an error.",
		Parameters: []function.Parameter{
			function.MapParameter{
				Name:                "entries",
				MarkdownDescription: "Entries to render",
				ElementType:         types.StringType,
			},
		},
		Return: function.StringReturn{},
	}
}

func (f encodeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var raw map[string]types.String

	resp.Error = req.Arguments.Get(ctx, &raw)
	if resp.Error != nil {
		return
	}

	entries := make(map[string]string, len(raw))

	for key, value := range raw {
		if value.IsNull() {
			resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("value of key %q must not be null", key))

			return
		}

		entries[key] = value.ValueString()
	}

	content, err := dotenv.Encode(entries)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())

		return
	}

	resp.Error = resp.Result.Set(ctx, content)
}
