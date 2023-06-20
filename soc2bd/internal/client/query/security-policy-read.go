package query

import (
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"
)

type ReadSecurityPolicy struct {
	SecurityPolicy *gqlSecurityPolicy `graphql:"securityPolicy(id: $id, name: $name)"`
}

func (q ReadSecurityPolicy) IsEmpty() bool {
	return q.SecurityPolicy == nil
}

type gqlSecurityPolicy struct {
	IDName
}

func (q ReadSecurityPolicy) ToModel() *model.SecurityPolicy {
	if q.SecurityPolicy == nil {
		return nil
	}

	return q.SecurityPolicy.ToModel()
}

func (q *gqlSecurityPolicy) ToModel() *model.SecurityPolicy {
	return &model.SecurityPolicy{
		ID:   string(q.ID),
		Name: q.Name,
	}
}
