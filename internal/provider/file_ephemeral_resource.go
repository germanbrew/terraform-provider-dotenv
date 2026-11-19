package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ ephemeral.EphemeralResource = &fileDotEnvEphemeralResource{}

// NewFileDotEnvEphemeralResource returns the `dotenv` ephemeral resource.
func NewFileDotEnvEphemeralResource() ephemeral.EphemeralResource {
	return &fileDotEnvEphemeralResource{}
}

// fileDotEnvEphemeralResource defines the ephemeral resource implementation.
type fileDotEnvEphemeralResource struct{}

func (d *fileDotEnvEphemeralResource) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName
}

func (d *fileDotEnvEphemeralResource) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Reads and provides all entries of one or more dotenv files.\n\n" + descSupportedFormats + "\n\n" +
			"-> If you only need a specific value you can use the " +
			"[`get`](" + docsBaseURL + "/functions/get) provider function.\n\n" +
			"Ephemeral resources are not stored in the state. ",

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
				MarkdownDescription: "Key-Value entries of the dotenv file(s). The values are considered sensitive.",
				Computed:            true,
				Sensitive:           true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *fileDotEnvEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data dotenvFileModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.load(ctx)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
