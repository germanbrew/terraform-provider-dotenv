package provider

import (
	"context"
	"fmt"

	"github.com/germanbrew/terraform-provider-dotenv/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ function.Function = GetByKeyFunction{}

func NewGetByKeyFunction() function.Function {
	return GetByKeyFunction{}
}

type GetByKeyFunction struct{}

func (r GetByKeyFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "get_by_key"
}

func (r GetByKeyFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Get by key function",
		MarkdownDescription: "Reads and provides a single entry of a .env file by its key\n\n" +
			"All supported formats can be found [here](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs#supported-formats).",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "key",
				MarkdownDescription: "Name of the key",
			},
			function.StringParameter{
				Name:                "filename",
				MarkdownDescription: "Path to the dotenv file",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r GetByKeyFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var (
		key      string
		filename string
	)

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &key, &filename))

	if resp.Error != nil {
		return
	}

	if filename == "" {
		filename = ".env"
		tflog.Info(ctx, "No file name specified, so the default is used: "+filename)
	}

	entries, err := utils.ParseDotEnvFile(filename)
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())

		return
	}

	value, exists := entries[key]
	if !exists {
		resp.Error = function.NewFuncError(fmt.Sprintf("Could not find key '%s' in file '%s'", key, filename))

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, value))
}
