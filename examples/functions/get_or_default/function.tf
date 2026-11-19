output "log_level" {
  value = provider::dotenv::get_or_default("LOG_LEVEL", "info")
}

output "port" {
  value = provider::dotenv::get_or_default("PORT", "8080", "${path.module}/app.env")
}
