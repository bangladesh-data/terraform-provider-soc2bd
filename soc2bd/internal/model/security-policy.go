package model

import "github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"

type SecurityPolicy struct {
	ID   string
	Name string
}

func (s SecurityPolicy) ToTerraform() interface{} {
	return map[string]interface{}{
		attr.ID:   s.ID,
		attr.Name: s.Name,
	}
}
