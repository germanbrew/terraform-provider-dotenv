package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func functionTest(t *testing.T, step resource.TestStep) {
	t.Helper()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{step},
	})
}

func TestGetFunction_KnownKey(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "test" {
  value = provider::dotenv::get("EXAMPLE_STRING", "./testdata/test.env")
}

output "addition" {
  value = provider::dotenv::get("EXAMPLE_INT", "./testdata/test.env") + 50
}
`,
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("Example v@lue!")),
			statecheck.ExpectKnownOutputValue("addition", knownvalue.Int64Exact(150)),
		},
	})
}

func TestGetFunction_DefaultFilename(t *testing.T) {
	// The default .env is resolved relative to the working directory of the
	// test, which is the provider package directory.
	t.Chdir("testdata")

	functionTest(t, resource.TestStep{
		Config: `
output "test" {
  value = provider::dotenv::get("SOME_VAR")
}
`,
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("from-default-file")),
		},
	})
}

func TestGetFunction_UnknownKey(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "test" {
  value = provider::dotenv::get("DOES_NOT_EXIST", "./testdata/test.env")
}
`,
		ExpectError: expectErrorText("Could not find key 'DOES_NOT_EXIST'"),
	})
}

func TestGetFunction_UnknownFile(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "test" {
  value = provider::dotenv::get("EXAMPLE_STRING", "./testdata/unknown.env")
}
`,
		ExpectError: regexp.MustCompile(ErrFileNotFound.Error()),
	})
}

func TestGetFunction_TooManyFilenames(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "test" {
  value = provider::dotenv::get("EXAMPLE_STRING", "./testdata/test.env", "./testdata/base.env")
}
`,
		ExpectError: regexp.MustCompile(`at most one filename can be given`),
	})
}

func TestGetByKeyFunction_Deprecated(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "test" {
  value = provider::dotenv::get_by_key("EXAMPLE_STRING", "./testdata/test.env")
}
`,
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("Example v@lue!")),
		},
	})
}

func TestHasFunction(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "exists" {
  value = provider::dotenv::has("EXAMPLE_STRING", "./testdata/test.env")
}

output "missing" {
  value = provider::dotenv::has("DOES_NOT_EXIST", "./testdata/test.env")
}
`,
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownOutputValue("exists", knownvalue.Bool(true)),
			statecheck.ExpectKnownOutputValue("missing", knownvalue.Bool(false)),
		},
	})
}

func TestHasFunction_UnknownFile(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "test" {
  value = provider::dotenv::has("EXAMPLE_STRING", "./testdata/unknown.env")
}
`,
		ExpectError: regexp.MustCompile(ErrFileNotFound.Error()),
	})
}

func TestGetOrDefaultFunction(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "exists" {
  value = provider::dotenv::get_or_default("EXAMPLE_STRING", "fallback", "./testdata/test.env")
}

output "missing" {
  value = provider::dotenv::get_or_default("DOES_NOT_EXIST", "fallback", "./testdata/test.env")
}
`,
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownOutputValue("exists", knownvalue.StringExact("Example v@lue!")),
			statecheck.ExpectKnownOutputValue("missing", knownvalue.StringExact("fallback")),
		},
	})
}

func TestDecodeFunction(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "decoded" {
  value = provider::dotenv::decode(file("./testdata/base.env"))
}

output "inline" {
  value = provider::dotenv::decode("A=1\nB=\"two words\" # comment\n")
}
`,
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownOutputValue("decoded", knownvalue.MapExact(map[string]knownvalue.Check{
				"NAME":      knownvalue.StringExact("base"),
				"ONLY_BASE": knownvalue.StringExact("1"),
				"URL":       knownvalue.StringExact("http://base.local"),
			})),
			statecheck.ExpectKnownOutputValue("inline", knownvalue.MapExact(map[string]knownvalue.Check{
				"A": knownvalue.StringExact("1"),
				"B": knownvalue.StringExact("two words"),
			})),
		},
	})
}

func TestDecodeFunction_Invalid(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "test" {
  value = provider::dotenv::decode("this is invalid")
}
`,
		ExpectError: expectErrorText(ErrInvalidLine.Error()),
	})
}

func TestEncodeFunction(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
locals {
  entries = {
    PLAIN     = "value"
    SPACES    = "hello world"
    MULTILINE = "line1\nline2"
    DOLLAR    = "$HOME"
  }
}

output "encoded" {
  value = provider::dotenv::encode(local.entries)
}

output "round_trip" {
  value = provider::dotenv::decode(provider::dotenv::encode(local.entries)) == tomap(local.entries)
}
`,
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownOutputValue("encoded", knownvalue.StringExact(
				"DOLLAR='$HOME'\nMULTILINE=\"line1\\nline2\"\nPLAIN=value\nSPACES='hello world'\n",
			)),
			statecheck.ExpectKnownOutputValue("round_trip", knownvalue.Bool(true)),
		},
	})
}

func TestEncodeFunction_InvalidKey(t *testing.T) {
	functionTest(t, resource.TestStep{
		Config: `
output "test" {
  value = provider::dotenv::encode({ "invalid key" = "value" })
}
`,
		ExpectError: regexp.MustCompile(`invalid key "invalid key"`),
	})
}
