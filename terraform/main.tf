terraform {
  required_providers {
    dotenv = {
      source  = "germanbrew/dotenv"
      version = "1.0.0"
    }
  }

  required_version = ">= 1.8.0"
}

provider "dotenv" {}

data "dotenv" "test" {}

output "int" {
  value = data.dotenv.test.entries.INT + 2000
}

output "string" {
  value = provider::dotenv::get_by_key("STRING", ".env")
}

output "bool" {
  value = provider::dotenv::get_by_key("BOOL", ".env") ? "is true" : " is false"
}

