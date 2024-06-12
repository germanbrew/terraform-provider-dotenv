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
