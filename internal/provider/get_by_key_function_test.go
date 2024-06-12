package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestGetByKeyFunction_Known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
			  		value = provider::dotenv::get_by_key("EXAMPLE_STRING", "./testdata/test.env")
				}

				output "addition" {
			  		value = provider::dotenv::get_by_key("EXAMPLE_INT", "./testdata/test.env") + 50
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("Example v@lue!")),
					statecheck.ExpectKnownOutputValue("addition", knownvalue.Int64Exact(150)),
				},
			},
		},
	})
}

func TestGetByKeyFunction_UnknownKey(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
			  		value = provider::dotenv::get_by_key("DOES_NOT_EXIST", "./testdata/test.env")
				}
				`,
				ExpectError: regexp.MustCompile(`Could not find key`),
			},
		},
	})
}

func TestGetByKeyFunction_UnknownFileName(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
			  		value = provider::dotenv::get_by_key("EXAMPLE_STRING", "./testdata/unknown.env")
				}
				`,
				ExpectError: regexp.MustCompile(`no such file or directory`),
			},
		},
	})
}
