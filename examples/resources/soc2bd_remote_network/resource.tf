provider "soc2bd" {
  api_token = "1234567890abcdef"
  network   = "mynetwork"
}

resource "soc2bd_remote_network" "aws_network" {
  name = "aws_remote_network"
}
