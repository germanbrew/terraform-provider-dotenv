provider "dotenv" {}

data "dotenv_file" "app" {
  filename  = "./application/app.env"
  sensitive = true
}

provider "some_provider" {
  token = data.dotenv_file.app.SECRET_KEY
}
