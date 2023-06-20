resource "soc2bd_remote_network" "aws_remote_network" {
  name = "aws-remote-network"
}

resource "soc2bd_connector" "aws_connector" {
  remote_network_id = soc2bd_remote_network.aws_remote_network.id
}

resource "soc2bd_connector_tokens" "aws_connector_tokens" {
  connector_id = soc2bd_connector.aws_connector.id
}

