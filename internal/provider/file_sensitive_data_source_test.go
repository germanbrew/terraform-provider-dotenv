package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSource_SensitiveDotEnvFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleSensitiveDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dotenv_sensitive.test", "entries.EXAMPLE_STRING", "Example v@lue!"),
					resource.TestCheckResourceAttr("data.dotenv_sensitive.test", "entries.EXAMPLE_INT", "100"),
					resource.TestCheckResourceAttr("data.dotenv_sensitive.test", "entries.EXAMPLE_FLOAT", "1.23"),
					resource.TestCheckResourceAttr("data.dotenv_sensitive.test", "entries.SOME_VAR", "someval"),
					resource.TestCheckResourceAttr("data.dotenv_sensitive.test", "entries.BAR", "BAZ"),
					resource.TestCheckResourceAttr("data.dotenv_sensitive.test", "entries.FOO", "BAR"),
					resource.TestCheckResourceAttr("data.dotenv_sensitive.test", "entries.YAML_FOO", "bar"),
				),
			},
		},
	})
}

func TestAccDataSource_SensitiveDotEnvFile_UnknownFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccUnknownSensitiveDataSourceConfig,
				ExpectError: regexp.MustCompile(fmt.Sprintf("%s: %s", "testdata/unknown.env", NotFoundError)),
			},
		},
	})
}

func TestAccDataSource_SensitiveDotEnvFile_UnknownKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleSensitiveDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.dotenv_sensitive.test", "entries.unknown", "invalid"),
				),
				ExpectError: regexp.MustCompile(`data.dotenv_sensitive.test: Attribute 'entries.unknown' not found`),
			},
		},
	})
}

const testAccExampleSensitiveDataSourceConfig = `
data "dotenv_sensitive" "test" {
  filename = "./testdata/test.env"
}
`

const testAccUnknownSensitiveDataSourceConfig = `
data "dotenv" "test" {
	  filename = "./testdata/unknown.env"
}
`
