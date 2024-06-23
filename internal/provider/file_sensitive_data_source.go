package provider

import (
	"context"
	"fmt"

	"github.com/germanbrew/terraform-provider-dotenv/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &sensitiveFileDotEnvDataSource{}

func NewSensitiveFileDotEnvDataSource() datasource.DataSource {
	return &sensitiveFileDotEnvDataSource{}
}

// sensitiveFileDotEnvDataSource defines the data source implementation.
type sensitiveFileDotEnvDataSource struct{}

// sensitiveFileDotEnvDataSourceModel describes the data source data model.
type sensitiveFileDotEnvDataSourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Entries  types.Map    `tfsdk:"entries"`
}

func (d *sensitiveFileDotEnvDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive"
}

func (d *sensitiveFileDotEnvDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Reads and provides all entries of a dotenv file as sensitive data.\n\n" +
			"All supported formats can be found [here](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs#supported-formats).\n\n" +
			"-> If you only need a specific value and/or do not want to store the contents of the file in the state, " +
			"you can use the [`get_by_key`](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs/functions/get_by_key) provider function.",

		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				MarkdownDescription: "`Default: .env` Path to the dotenv file",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"entries": schema.MapAttribute{
				MarkdownDescription: "Key-Value entries of the dotenv file. The values are by default considered sensitive. " +
					"If you want to treat the values as non-sensitive, you can use the " +
					"[`dotenv`](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs/data-sources/dotenv) data source " +
					"or [`nonsensitive()`](https://developer.hashicorp.com/terraform/language/functions/nonsensitive) function.",
				Computed:    true,
				Sensitive:   true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *sensitiveFileDotEnvDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *sensitiveFileDotEnvDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data sensitiveFileDotEnvDataSourceModel

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

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
