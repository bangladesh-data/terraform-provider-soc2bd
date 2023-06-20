package query

import (
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"
)

type CreateServiceAccountKey struct {
	ServiceAccountKeyEntityCreateResponse `graphql:"serviceAccountKeyCreate(expirationTime: $expirationTime, serviceAccountId: $serviceAccountId, name: $name)"`
}

func (q CreateServiceAccountKey) IsEmpty() bool {
	return q.Entity == nil
}

type ServiceAccountKeyEntityCreateResponse struct {
	ServiceAccountKeyEntityResponse
	Token string
}

func (q CreateServiceAccountKey) ToModel() (*model.ServiceKey, error) {
	if q.Entity == nil {
		return nil, nil //nolint
	}

	serviceKey, err := q.Entity.ToModel()
	if err != nil {
		return nil, err
	}

	serviceKey.Token = q.Token

	return serviceKey, nil
}
