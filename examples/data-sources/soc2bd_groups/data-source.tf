data "soc2bd_groups" "foo" {
  name = "<your group's name>"
}

# Group names are not constrained to be unique within Soc2bd,
# so it is possible that this data source will return multiple list items.
