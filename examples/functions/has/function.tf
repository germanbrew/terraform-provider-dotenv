locals {
  debug = provider::dotenv::has("DEBUG") ? provider::dotenv::get("DEBUG") == "true" : false
}
