provider "soc2bd" {
  api_token = "1234567890abcdef"
  network   = "mynetwork"
}

resource "soc2bd_service_account" "github_actions_prod" {
  name = "Github Actions PROD"
}

resource "soc2bd_service_account_key" "github_key" {
  name = "Github Actions PROD key"
  service_account_id = soc2bd_service_account.github_actions_prod.id
}
