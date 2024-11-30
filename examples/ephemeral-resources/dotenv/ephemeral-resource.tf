ephemeral "dotenv" "database" {
  filename = "./database.env"
}

provider "postgresql" {
  host     = data.aws_db_instance.example.address
  port     = data.aws_db_instance.example.port
  username = ephemeral.dotenv.database.entries.USERNAME
  password = ephemeral.dotenv.database.entries.PASSWORD
}
