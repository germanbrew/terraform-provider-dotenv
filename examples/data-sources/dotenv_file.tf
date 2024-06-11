# Specify specifc dotfile path
data "dotenv_file" "database_env" {
  filename = "./database.env"
}

# Use .env file in current directory
data "dotenv_file" "local_env" {}
