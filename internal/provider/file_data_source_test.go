package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSource_DotEnvFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.EXAMPLE_STRING", "Example v@lue!"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.EXAMPLE_INT", "100"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.EXAMPLE_FLOAT", "1.23"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.SOME_VAR", "someval"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.BAR", "BAZ"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.FOO", "BAR"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.YAML_FOO", "bar"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
data "dotenv" "test" {
  filename = "./testdata/test.env"
}
`
