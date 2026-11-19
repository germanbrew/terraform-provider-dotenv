# Read a key from the .env file in the current directory
output "api_url" {
  value = provider::dotenv::get("API_URL")
}

# Read a key from a specific file, relative to the module directory
resource "some_resource" "example" {
  value = provider::dotenv::get("SECRET_KEY", "${path.module}/app.env")
}

# Values are strings, but can be converted implicitly
output "addition" {
  value = provider::dotenv::get("EXAMPLE_INT", "./testdata/test.env") + 50
}
