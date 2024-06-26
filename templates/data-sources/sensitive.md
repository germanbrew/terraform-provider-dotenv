---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dotenv_sensitive Data Source - dotenv"
subcategory: ""
description: |-
  Reads and provides all entries of a dotenv file as sensitive data.
  All supported formats can be found here https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs#supported-formats.
  -> If you only need a specific value and/or do not want to store the contents of the file in the state, you can use the get_by_key https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs/functions/get_by_key provider function.
---

# dotenv_sensitive (Data Source)

Reads and provides all entries of a dotenv file as sensitive data.

All supported formats can be found [here](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs#supported-formats).

-> If you only need a specific value and/or do not want to store the contents of the file in the state, you can use the [`get_by_key`](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs/functions/get_by_key) provider function.

## Example Usage

```terraform
# Use the default .env in current directory
data "dotenv_sensitive" "local" {}

output "local" {
  value = data.dotenv_sensitive.local.entries.EXAMPLE_KEY
}

# With specific .env file path
data "dotenv_sensitive" "database" {
  filename = "./database.env"
}

output "foo" {
  value     = data.dotenv_sensitive.database.entries.PASSWORD
  sensitive = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filename` (String) `Default: .env` Path to the dotenv file

### Read-Only

- `entries` (Map of String, Sensitive) Key-Value entries of the dotenv file. The values are by default considered sensitive. If you want to treat the values as non-sensitive, you can use the [`dotenv`](https://registry.terraform.io/providers/germanbrew/dotenv/latest/docs/data-sources/dotenv) data source or [`nonsensitive()`](https://developer.hashicorp.com/terraform/language/functions/nonsensitive) function.
