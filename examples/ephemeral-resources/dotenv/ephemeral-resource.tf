ephemeral "dotenv" "database" {
  filename   = "${path.module}/database.env"
  expand_env = false
}

provider "postgresql" {
  host     = data.aws_db_instance.example.address
  port     = data.aws_db_instance.example.port
  username = ephemeral.dotenv.database.entries.USERNAME
  password = ephemeral.dotenv.database.entries.PASSWORD
}
