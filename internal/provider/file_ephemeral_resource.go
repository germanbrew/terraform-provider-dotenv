package provider

import (
	"context"
	"fmt"

	"github.com/germanbrew/terraform-provider-dotenv/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ ephemeral.EphemeralResource = &fileDotEnvEphemeralResource{}

func NewFileDotEnvEphemeralResource() ephemeral.EphemeralResource {
	return &fileDotEnvEphemeralResource{}
}

// fileDotEnvEphemeralResource defines the data source implementation.
type fileDotEnvEphemeralResource struct{}

// fileDotEnvEphemeralResourceModel describes the data source data model.
type fileDotEnvEphemeralResourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Entries  types.Map    `tfsdk:"entries"`
}

func (d *fileDotEnvEphemeralResource) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName // + "_file"
}

func (d *fileDotEnvEphemeralResource) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Reads and provides all entries of a dotenv file.\n\n" +
			"All supported formats can be found [here](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs#supported-formats).\n\n" +
			"-> If you only need a specific value you can use the " +
			"[`get_by_key`](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs/functions/get_by_key) provider function.\n\n" +
			"Ephemeral resources are not stored in the state. ",

		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				MarkdownDescription: "`Default: .env` Path to the dotenv file",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"entries": schema.MapAttribute{
				MarkdownDescription: "Key-Value entries of the dotenv file.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *fileDotEnvEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *fileDotEnvEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data fileDotEnvEphemeralResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	filename := data.Filename.ValueString()

	if filename == "" {
		filename = ".env"
		tflog.Info(ctx, "No file name specified, so the default is used: "+filename)
	}

	parsedEntries, err := utils.ParseDotEnvFile(filename)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Parsing contents of file %s failed: %s", filename, err))

		return
	}

	entries := make(map[string]attr.Value, len(parsedEntries))
	for key, value := range parsedEntries {
		entries[key] = types.StringValue(value)

		if err != nil {
			resp.Diagnostics.AddError("Conversion Error", fmt.Sprintf("Failed to convert key %s value %s: %s", key, value, err))

			return
		}
	}

	tflog.Debug(ctx, "Parsing the file was successful")

	data.Filename = types.StringValue(filename)
	data.Entries, _ = types.MapValue(types.StringType, entries)

	// Save data into ephemeral result data
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
