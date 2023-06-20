provider "soc2bd" {
  api_token = "1234567890abcdef"
  network   = "mynetwork"
}

resource "soc2bd_service_account" "github_actions_prod" {
  name = "Github Actions PROD"
}
