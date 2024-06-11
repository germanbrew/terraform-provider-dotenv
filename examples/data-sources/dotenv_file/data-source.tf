# With specifc dotfile path and without storing the contents to the TF state
data "dotenv_file" "database" {
  filename   = "./database.env"
  keep_local = true
}

output "foo" {
  value = data.dotenv_file.database.PASSWORD
}

# Use the default .env in current directory
data "dotenv_file" "local" {}

output "local" {
  value = data.dotenv_file.local.EXAMPLE_KEY
}
