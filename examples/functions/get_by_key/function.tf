# Deprecated: use provider::dotenv::get("SECRET_KEY", ".env") instead
resource "some_resource" "example" {
  value = provider::dotenv::get_by_key("SECRET_KEY", ".env")
}
