# With specifc dotfile path and without storing the contents to the TF state
data "dotenv" "database" {
  filename   = "./database.env"
  keep_local = true
}

output "foo" {
  value = data.dotenv.database.PASSWORD
}

# Use the default .env in current directory
data "dotenv" "local" {}

output "local" {
  value = data.dotenv.local.EXAMPLE_KEY
}
