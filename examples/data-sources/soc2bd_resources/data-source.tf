data "soc2bd_resources" "foo" {
  name = "<your resource's name>"
}

# Resource names are not constrained to be unique within Soc2bd,
# so it is possible that this data source will return multiple list items.
