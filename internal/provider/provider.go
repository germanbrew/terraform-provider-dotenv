package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure DotenvProvider satisfies various provider interfaces.
var (
	_ provider.Provider              = &DotenvProvider{}
	_ provider.ProviderWithFunctions = &DotenvProvider{}
)

// DotenvProvider defines the provider implementation.
type DotenvProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DotenvProviderModel describes the provider data model.
type DotenvProviderModel struct{}

func (p *DotenvProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dotenv"
	resp.Version = p.version
}

func (p *DotenvProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A utility Terraform provider for .env files.  \n" +
			"## Supported formats\n" +
			"This provider supports the following formats:  \n" +
			"```sh\n" +
			"SOME_VAR=someval\n" +
			"export BAR=BAZ\n" +
			"# I am a comment and that is OK\n" +
			"FOO=BAR # comments at line end are OK too\n" +
			"```\n" +
			"You can also do a YAML(ish) style:\n" +
			"```sh  \n" +
			"FOO: bar\n" +
			"BAR: baz\n" +
			"```",
	}
}

func (p *DotenvProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data DotenvProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *DotenvProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *DotenvProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFileDotEnvDataSource,
	}
}

func (p *DotenvProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewGetByKeyFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DotenvProvider{
			version: version,
		}
	}
}
