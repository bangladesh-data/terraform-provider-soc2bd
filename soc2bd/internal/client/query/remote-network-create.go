package query

import "github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"

type CreateRemoteNetwork struct {
	RemoteNetworkEntityResponse `graphql:"remoteNetworkCreate(name: $name, isActive: $isActive, location: $location)"`
}

func (q CreateRemoteNetwork) IsEmpty() bool {
	return q.Entity == nil
}

type RemoteNetworkEntityResponse struct {
	Entity *gqlRemoteNetwork
	OkError
}

func (q CreateRemoteNetwork) ToModel() *model.RemoteNetwork {
	if q.Entity == nil {
		return nil
	}

	return q.Entity.ToModel()
}
