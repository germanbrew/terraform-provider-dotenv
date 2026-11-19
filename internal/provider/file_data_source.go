package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &fileDotEnvDataSource{}

// NewFileDotEnvDataSource returns the `dotenv` data source.
func NewFileDotEnvDataSource() datasource.DataSource {
	return &fileDotEnvDataSource{
		typeNameSuffix: "",
		sensitive:      false,
		description:    "Reads and provides all entries of one or more dotenv files.",
		entriesDescription: "Key-Value entries of the dotenv file(s). The values are by default considered nonsensitive. " +
			"If you want to treat the values as sensitive, you can use the " +
			"[`dotenv_sensitive`](" + docsBaseURL + "/data-sources/sensitive) data source " +
			"or [`sensitive()`](https://developer.hashicorp.com/terraform/language/functions/sensitive) function.",
	}
}

// NewSensitiveFileDotEnvDataSource returns the `dotenv_sensitive` data source.
func NewSensitiveFileDotEnvDataSource() datasource.DataSource {
	return &fileDotEnvDataSource{
		typeNameSuffix: "_sensitive",
		sensitive:      true,
		description:    "Reads and provides all entries of one or more dotenv files as sensitive data.",
		entriesDescription: "Key-Value entries of the dotenv file(s). The values are by default considered sensitive. " +
			"If you want to treat the values as non-sensitive, you can use the " +
			"[`dotenv`](" + docsBaseURL + "/data-sources/dotenv) data source " +
			"or [`nonsensitive()`](https://developer.hashicorp.com/terraform/language/functions/nonsensitive) function.",
	}
}

// fileDotEnvDataSource implements the `dotenv` and `dotenv_sensitive` data
// sources, which only differ in the sensitivity of their entries.
type fileDotEnvDataSource struct {
	typeNameSuffix     string
	sensitive          bool
	description        string
	entriesDescription string
}

func (d *fileDotEnvDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.typeNameSuffix
}

func (d *fileDotEnvDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: d.description + "\n\n" + descSupportedFormats + "\n\n" +
			"-> If you only need a specific value and/or do not want to store the contents of the file in the state, " +
			"you can use the [`get`](" + docsBaseURL + "/functions/get) provider function " +
			"or the [`dotenv`](" + docsBaseURL + "/ephemeral-resources/dotenv) ephemeral resource.",

		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				MarkdownDescription: descFilename,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ConflictsWith(path.MatchRoot("filenames")),
				},
			},
			"filenames": schema.ListAttribute{
				MarkdownDescription: descFilenames,
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
			},
			"allow_missing": schema.BoolAttribute{
				MarkdownDescription: descAllowMissing,
				Optional:            true,
			},
			"expand_env": schema.BoolAttribute{
				MarkdownDescription: descExpandEnv,
				Optional:            true,
			},
			"prefix": schema.StringAttribute{
				MarkdownDescription: descPrefix,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"trim_prefix": schema.BoolAttribute{
				MarkdownDescription: descTrimPrefix,
				Optional:            true,
				Validators: []validator.Bool{
					boolvalidator.AlsoRequires(path.MatchRoot("prefix")),
				},
			},
			"entries": schema.MapAttribute{
				MarkdownDescription: d.entriesDescription,
				Computed:            true,
				Sensitive:           d.sensitive,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *fileDotEnvDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dotenvFileModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.load(ctx)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
