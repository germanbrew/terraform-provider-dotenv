---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dotenv Data Source - dotenv"
subcategory: ""
description: |-
  Reads and provides all entries of a dotenv file.
  All supported formats can be found here https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs#supported-formats.
  -> If you only need a specific value and/or do not want to store the contents of the file in the state, you can use the get_by_key https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs/functions/get_by_key provider function.
---

# dotenv (Data Source)

Reads and provides all entries of a dotenv file.

All supported formats can be found [here](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs#supported-formats).

-> If you only need a specific value and/or do not want to store the contents of the file in the state, you can use the [`get_by_key`](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs/functions/get_by_key) provider function.

## Example Usage

```terraform
# Use the default .env in current directory
data "dotenv" "local" {}

output "local" {
  value = data.dotenv.local.entries.EXAMPLE_KEY
}

# With specific .env file path
data "dotenv" "database" {
  filename = "./database.env"
}

output "foo" {
  value = data.dotenv.database.entries.PASSWORD
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filename` (String) `Default: .env` Path to the dotenv file

### Read-Only

- `entries` (Map of String) Key-Value entries of the dotenv file. The values are by default considered nonsensitive. If you want to treat the values as sensitive, you can use the [`dotenv_sensitive`](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs/data-sources/sensitive) data source or [`sensitive()`](https://developer.hashicorp.com/terraform/language/functions/sensitive) function.
