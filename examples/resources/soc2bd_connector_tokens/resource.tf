provider "soc2bd" {
  api_token = "1234567890abcdef"
  network   = "mynetwork"
}

resource "soc2bd_remote_network" "aws_network" {
  name = "aws_remote_network"
}

resource "soc2bd_connector" "aws_connector" {
  remote_network_id = soc2bd_remote_network.aws_network.id
}

resource "soc2bd_connector_tokens" "aws_connector_tokens" {
  connector_id = soc2bd_connector.aws_connector.id
}
