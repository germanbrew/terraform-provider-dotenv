# Use the default .env in current directory
data "dotenv_sensitive" "local" {}

output "local" {
  value     = data.dotenv_sensitive.local.entries.EXAMPLE_KEY
  sensitive = true
}

# With specific .env file path, relative to the module directory
data "dotenv_sensitive" "database" {
  filename   = "${path.module}/database.env"
  expand_env = false
}

output "foo" {
  value     = data.dotenv_sensitive.database.entries.PASSWORD
  sensitive = true
}
