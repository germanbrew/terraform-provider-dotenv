// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
					resource.TestCheckResourceAttr("data.dotenv_file.test", "entries.EXAMPLE_STRING", "Example v@lue!"),
					resource.TestCheckResourceAttr("data.dotenv_file.test", "entries.EXAMPLE_INT", "345"),
					resource.TestCheckResourceAttr("data.dotenv_file.test", "entries.EXAMPLE_FLOAT", "1.23"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
data "dotenv_file" "test" {
  filename = "./testdata/test.env"
}
`
