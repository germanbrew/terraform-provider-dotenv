package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSource_DotEnvFile_ExpandEnv(t *testing.T) {
	t.Setenv("DOTENV_TEST_HOST_VAR", "from-host")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "dotenv" "default" {
  filename = "./testdata/expand.env"
}

data "dotenv" "no_env" {
  filename   = "./testdata/expand.env"
  expand_env = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// default: host environment is used (backwards compatible)
					resource.TestCheckResourceAttr("data.dotenv.default", "entries.FROM_FILE", "local-value"),
					resource.TestCheckResourceAttr("data.dotenv.default", "entries.FROM_HOST", "from-host"),
					resource.TestCheckResourceAttr("data.dotenv.default", "entries.WITH_DEFAULT", "fallback"),
					resource.TestCheckResourceAttr("data.dotenv.default", "entries.ESCAPED", "$LOCAL"),
					resource.TestCheckResourceAttr("data.dotenv.default", "entries.SINGLE_QUOTED", "${LOCAL}"),
					// expand_env = false: only references within the file are resolved
					resource.TestCheckResourceAttr("data.dotenv.no_env", "entries.FROM_FILE", "local-value"),
					resource.TestCheckResourceAttr("data.dotenv.no_env", "entries.FROM_HOST", ""),
					resource.TestCheckResourceAttr("data.dotenv.no_env", "entries.WITH_DEFAULT", "fallback"),
				),
			},
		},
	})
}

func TestAccDataSource_DotEnvFile_Filenames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "dotenv" "test" {
  filenames = ["./testdata/base.env", "./testdata/override.env"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.%", "5"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.NAME", "override"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.ONLY_BASE", "1"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.ONLY_OVERRIDE", "2"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.FROM_BASE", "1"),
					// expanded while reading base.env, before it was overridden
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.URL", "http://base.local"),
					resource.TestCheckNoResourceAttr("data.dotenv.test", "filename"),
				),
			},
		},
	})
}

func TestAccDataSource_DotEnvFile_FilenameConflict(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "dotenv" "test" {
  filename  = "./testdata/base.env"
  filenames = ["./testdata/override.env"]
}
`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestAccDataSource_DotEnvFile_AllowMissing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "dotenv" "partial" {
  filenames     = ["./testdata/base.env", "./testdata/does-not-exist.env"]
  allow_missing = true
}

data "dotenv" "none" {
  filename      = "./testdata/does-not-exist.env"
  allow_missing = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dotenv.partial", "entries.NAME", "base"),
					resource.TestCheckResourceAttr("data.dotenv.none", "entries.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSource_DotEnvFile_AllowMissingInvalidStillFails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "dotenv" "test" {
  filename      = "./testdata/invalid.env"
  allow_missing = true
}
`,
				ExpectError: expectErrorText(ErrInvalidLine.Error()),
			},
		},
	})
}

func TestAccDataSource_DotEnvFile_Prefix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "dotenv" "filtered" {
  filename = "./testdata/prefixed.env"
  prefix   = "APP_"
}

data "dotenv" "trimmed" {
  filename    = "./testdata/prefixed.env"
  prefix      = "APP_"
  trim_prefix = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dotenv.filtered", "entries.%", "3"),
					resource.TestCheckResourceAttr("data.dotenv.filtered", "entries.APP_NAME", "demo"),
					resource.TestCheckNoResourceAttr("data.dotenv.filtered", "entries.DB_HOST"),
					// APP_ would become an empty key and is skipped
					resource.TestCheckResourceAttr("data.dotenv.trimmed", "entries.%", "2"),
					resource.TestCheckResourceAttr("data.dotenv.trimmed", "entries.NAME", "demo"),
					resource.TestCheckResourceAttr("data.dotenv.trimmed", "entries.PORT", "8080"),
				),
			},
		},
	})
}

func TestAccDataSource_DotEnvFile_TrimPrefixRequiresPrefix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "dotenv" "test" {
  filename    = "./testdata/prefixed.env"
  trim_prefix = true
}
`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestAccDataSource_DotEnvFile_Multiline(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "dotenv" "test" {
  filename = "./testdata/multiline.env"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.MULTILINE", "line1\nline2"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.ESCAPED_NEWLINE", "a\nb"),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.JSON", `{"a":1}`),
					resource.TestCheckResourceAttr("data.dotenv.test", "entries.QUOTED", `say "hi"`),
				),
			},
		},
	})
}

func TestAccDataSource_DotEnvSensitive_Options(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "dotenv_sensitive" "test" {
  filenames   = ["./testdata/prefixed.env", "./testdata/missing.env"]
  prefix      = "DB_"
  trim_prefix = true
  allow_missing = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dotenv_sensitive.test", "entries.%", "1"),
					resource.TestCheckResourceAttr("data.dotenv_sensitive.test", "entries.HOST", "localhost"),
				),
			},
		},
	})
}
