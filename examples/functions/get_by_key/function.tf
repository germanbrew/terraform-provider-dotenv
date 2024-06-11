resource "some_ressource" "example" {
  value = provider::dotenv::get_by_key("SECRET_KEY", ".env")
}
