provider "dotenv" {}

data "dotenv" "app" {
  filename = "./application/app.env"
}

provider "some_provider" {
  token = data.dotenv.app.entries.SECRET_KEY
}
