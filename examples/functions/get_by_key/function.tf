resource "some_ressource" "example" {
  value = provider::dotenv::get_by_key("SECRET_KEY", ".env")
}

output "addition" {
  value = provider::dotenv::get_by_key("EXAMPLE_INT", "./testdata/test.env") + 50
}
