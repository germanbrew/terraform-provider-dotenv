package provider

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/germanbrew/terraform-provider-dotenv/internal/dotenv"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	defaultFilename = ".env"
	docsBaseURL     = "https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs"

	descSupportedFormats = "All supported formats can be found [here](" + docsBaseURL + "#supported-formats)."
	descRelativePaths    = "Relative paths are resolved against the directory Terraform is run in, not the module directory. " +
		"Use `${path.module}` to reference files next to the module, e.g. `\"${path.module}/.env\"`."

	descFilename  = "`Default: .env` Path to the dotenv file. " + descRelativePaths + " Conflicts with `filenames`."
	descFilenames = "List of dotenv files that are read in order and merged. Entries of later files override entries " +
		"of earlier files and can reference them in variable expansions. Conflicts with `filename`."
	descAllowMissing = "`Default: false` If `true`, files that do not exist are skipped instead of causing an error. " +
		"If no file exists, `entries` is an empty map."
	descExpandEnv = "`Default: true` Whether variable references (`$VAR`, `${VAR}`) that are not defined in the " +
		"dotenv file(s) are resolved from the environment of the Terraform process. If `false`, they resolve to an " +
		"empty string or their default (`${VAR:-default}`). Set this to `false` to get the same result on every " +
		"machine and to prevent environment variables of the host (e.g. CI secrets) from ending up in the values."
	descPrefix     = "Only include entries whose key starts with this prefix."
	descTrimPrefix = "`Default: false` Remove `prefix` from the keys of the returned entries. Requires `prefix`."

	descFunctionExpansion = "Variable references (`$VAR`, `${VAR}`) that are not defined in the file are resolved from " +
		"the environment of the Terraform process. Use the `dotenv` data source with `expand_env = false` if you need " +
		"results that do not depend on that environment."
)

// dotenvFileModel is the configuration and result model shared by the
// `dotenv` and `dotenv_sensitive` data sources and the `dotenv` ephemeral
// resource.
type dotenvFileModel struct {
	Filename     types.String `tfsdk:"filename"`
	Filenames    types.List   `tfsdk:"filenames"`
	AllowMissing types.Bool   `tfsdk:"allow_missing"`
	ExpandEnv    types.Bool   `tfsdk:"expand_env"`
	Prefix       types.String `tfsdk:"prefix"`
	TrimPrefix   types.Bool   `tfsdk:"trim_prefix"`
	Entries      types.Map    `tfsdk:"entries"`
}

// loadOptions describe which dotenv files are read and how.
type loadOptions struct {
	Filenames    []string
	AllowMissing bool
	ExpandEnv    bool
	Prefix       string
	TrimPrefix   bool
}

// load reads the files configured in the model and sets Entries (and the
// default Filename, if no file was configured).
func (m *dotenvFileModel) load(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics

	opts := loadOptions{
		AllowMissing: m.AllowMissing.ValueBool(),
		ExpandEnv:    m.ExpandEnv.IsNull() || m.ExpandEnv.ValueBool(),
		Prefix:       m.Prefix.ValueString(),
		TrimPrefix:   m.TrimPrefix.ValueBool(),
	}

	switch {
	case !m.Filenames.IsNull():
		diags.Append(m.Filenames.ElementsAs(ctx, &opts.Filenames, false)...)

		if diags.HasError() {
			return diags
		}
	case m.Filename.ValueString() != "":
		opts.Filenames = []string{m.Filename.ValueString()}
	default:
		tflog.Info(ctx, "No file name specified, so the default is used: "+defaultFilename)

		m.Filename = types.StringValue(defaultFilename)
		opts.Filenames = []string{defaultFilename}
	}

	entries, err := loadEntries(ctx, opts)
	if err != nil {
		diags.AddError("Parse Error", err.Error())

		return diags
	}

	m.Entries, diags = types.MapValueFrom(ctx, types.StringType, entries)

	return diags
}

// loadEntries reads and merges the given dotenv files.
func loadEntries(ctx context.Context, opts loadOptions) (map[string]string, error) {
	parseOpts := dotenv.Options{Base: map[string]string{}}
	if opts.ExpandEnv {
		parseOpts.LookupEnv = os.LookupEnv
	}

	for _, filename := range opts.Filenames {
		contents, err := os.ReadFile(filename)
		if err != nil {
			if opts.AllowMissing && errors.Is(err, fs.ErrNotExist) {
				tflog.Info(ctx, "Skipping missing dotenv file", map[string]any{"filename": filename})

				continue
			}

			return nil, fmt.Errorf("parsing contents of file %s failed: failed to read dotenv file: %w", filename, err)
		}

		parsed, err := dotenv.ParseWithOptions(string(contents), parseOpts)
		if err != nil {
			return nil, fmt.Errorf("parsing contents of file %s failed: %w", filename, err)
		}

		tflog.Debug(ctx, "Parsing the file was successful", map[string]any{"filename": filename})

		parseOpts.Base = parsed
	}

	return filterEntries(parseOpts.Base, opts.Prefix, opts.TrimPrefix), nil
}

func filterEntries(entries map[string]string, prefix string, trimPrefix bool) map[string]string {
	if prefix == "" {
		return entries
	}

	filtered := make(map[string]string, len(entries))

	for key, value := range entries {
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		if trimPrefix {
			key = strings.TrimPrefix(key, prefix)
			if key == "" {
				continue
			}
		}

		filtered[key] = value
	}

	return filtered
}

// readFileEntries reads a single dotenv file the way the provider functions
// do, i.e. with variable expansion from the process environment.
func readFileEntries(ctx context.Context, filename string) (map[string]string, error) {
	return loadEntries(ctx, loadOptions{Filenames: []string{filename}, ExpandEnv: true})
}
