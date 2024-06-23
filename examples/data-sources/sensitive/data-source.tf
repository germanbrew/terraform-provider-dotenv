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
