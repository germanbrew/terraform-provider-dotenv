package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccEphemeralResource_DotEnvFile_KnownKey(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleEphemeralResourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.env", tfjsonpath.New("data").AtMapKey("EXAMPLE_STRING"), knownvalue.StringExact("Example v@lue!")),
					statecheck.ExpectKnownValue("echo.env", tfjsonpath.New("data").AtMapKey("EXAMPLE_INT"), knownvalue.StringExact("100")),
					statecheck.ExpectKnownValue("echo.env", tfjsonpath.New("data").AtMapKey("EXAMPLE_FLOAT"), knownvalue.StringExact("1.23")),
					statecheck.ExpectKnownValue("echo.env", tfjsonpath.New("data").AtMapKey("SOME_VAR"), knownvalue.StringExact("someval")),
					statecheck.ExpectKnownValue("echo.env", tfjsonpath.New("data").AtMapKey("BAR"), knownvalue.StringExact("BAZ")),
					statecheck.ExpectKnownValue("echo.env", tfjsonpath.New("data").AtMapKey("FOO"), knownvalue.StringExact("BAR")),
					statecheck.ExpectKnownValue("echo.env", tfjsonpath.New("data").AtMapKey("YAML_FOO"), knownvalue.StringExact("bar")),
				},
			},
		},
	})
}

func TestAccEphemeralResource_DotEnvFile_UnknownKey(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleEphemeralResourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.env", tfjsonpath.New("data").AtMapKey("unknown"), knownvalue.StringExact("invalid")),
				},
				ExpectError: regexp.MustCompile(`path not found: specified key unknown not found in map at data.unknown`),
			},
		},
	})
}

func TestAccEphemeralResource_DotEnvFile_UnknownFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccUnknownEphemeralResourceConfig,
				ExpectError: regexp.MustCompile(fmt.Sprintf("%s: %s", "testdata/unknown.env", ErrFileNotFound)),
			},
		},
	})
}

func TestAccEphemeralResource_DotEnvFile_InvalidLine(t *testing.T) {
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      testAccInvalidEphemeralResourceConfig,
				ExpectError: regexp.MustCompile(fmt.Sprintf("%s: %s", ErrInvalidLine, "this is invalid")),
			},
		},
	})
}

// lintignore:AT004
const testAccExampleEphemeralResourceConfig = `
ephemeral "dotenv" "test" {
	filename = "./testdata/test.env"
}

provider "echo" {
	data = ephemeral.dotenv.test.entries
}

resource "echo" "env" {}
`

// lintignore:AT004
const testAccUnknownEphemeralResourceConfig = `
ephemeral "dotenv" "test" {
	filename = "./testdata/unknown.env"
}

provider "echo" {
	data = ephemeral.dotenv.test.entries
}

resource "echo" "env" {}
`

// lintignore:AT004
const testAccInvalidEphemeralResourceConfig = `
ephemeral "dotenv" "test" {
	filename = "./testdata/invalid.env"
}

provider "echo" {
	data = ephemeral.dotenv.test.entries
}

resource "echo" "env" {}
`
