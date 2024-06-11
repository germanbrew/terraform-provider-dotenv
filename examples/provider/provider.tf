provider "dotenv" {}

data "dotenv_file" "app" {
  filename = "./application/app.env"
}

provider "some_provider" {
  token = data.dotenv_file.app.SECRET_KEY
}
