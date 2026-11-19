# Generate a .env file, e.g. for docker compose
resource "local_file" "env" {
  filename = "${path.module}/generated.env"
  content = provider::dotenv::encode({
    DATABASE_HOST = aws_db_instance.example.address
    DATABASE_PORT = tostring(aws_db_instance.example.port)
    GREETING      = "Hello World"
  })
}

# Result:
# DATABASE_HOST=example.abc123.eu-central-1.rds.amazonaws.com
# DATABASE_PORT=5432
# GREETING='Hello World'
