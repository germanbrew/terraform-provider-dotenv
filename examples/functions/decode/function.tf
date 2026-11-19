# Parse the content of a file
locals {
  app = provider::dotenv::decode(file("${path.module}/app.env"))
}

# Parse a rendered template
locals {
  rendered = provider::dotenv::decode(templatefile("${path.module}/app.env.tftpl", {
    environment = "production"
  }))
}

output "database_url" {
  value = local.app["DATABASE_URL"]
}
