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
var _ datasource.DataSource = &fileDotEnvDataSource{}

func NewFileDotEnvDataSource() datasource.DataSource {
	return &fileDotEnvDataSource{}
}

// fileDotEnvDataSource defines the data source implementation.
type fileDotEnvDataSource struct{}

// fileDotEnvDataSourceModel describes the data source data model.
type fileDotEnvDataSourceModel struct {
	Filename  types.String `tfsdk:"filename"`
	KeepLocal types.Bool   `tfsdk:"keep_local"`
	Entries   types.Map    `tfsdk:"entries"`
}

func (d *fileDotEnvDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *fileDotEnvDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Reads and provides the entries of a dotenv file",

		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				MarkdownDescription: "`Default: .env` Path to the dotenv file",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"keep_local": schema.BoolAttribute{
				MarkdownDescription: "`Default: false` Set to `true` if the contents of the dotfile should not be " +
					"saved in the Terraform state and instead should always be loaded locally at runtime.",
				Optional: true,
			},
			"entries": schema.MapAttribute{
				MarkdownDescription: "Key-Value entries of the dotenv file",
				Computed:            true,
				ElementType:         types.StringType,
				Sensitive:           true,
			},
		},
	}
}

func (d *fileDotEnvDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *fileDotEnvDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data fileDotEnvDataSourceModel

	if !data.KeepLocal.ValueBool() {
		// Read Terraform configuration data into the model
		resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

		if resp.Diagnostics.HasError() {
			return
		}
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

	entries := make(map[string]attr.Value)
	for key, value := range parsedEntries {
		entries[key] = types.StringValue(value)
	}

	tflog.Debug(ctx, "Parsing the file was successful")

	data.Filename = types.StringValue(filename)
	data.Entries, _ = types.MapValue(types.StringType, entries)

	if !data.KeepLocal.ValueBool() {
		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}