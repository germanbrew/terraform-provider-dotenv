# Use the default .env in current directory
data "dotenv" "local" {}

output "local" {
  value = data.dotenv.local.entries.EXAMPLE_KEY
}

# With specific .env file path, relative to the module directory
data "dotenv" "database" {
  filename = "${path.module}/database.env"
}

output "foo" {
  value = data.dotenv.database.entries.PASSWORD
}

# Merge multiple files, later files override earlier ones.
# Files that do not exist are skipped.
data "dotenv" "merged" {
  filenames     = [".env", ".env.local"]
  allow_missing = true
}

# Only use entries starting with APP_ and remove the prefix from the keys,
# e.g. APP_NAME becomes NAME.
# Variable references are only resolved from the file itself, not from the
# environment Terraform runs in.
data "dotenv" "app" {
  filename    = ".env"
  prefix      = "APP_"
  trim_prefix = true
  expand_env  = false
}
