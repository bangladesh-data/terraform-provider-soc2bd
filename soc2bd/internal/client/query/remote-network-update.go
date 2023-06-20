package query

import "github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"

type UpdateRemoteNetwork struct {
	RemoteNetworkEntityResponse `graphql:"remoteNetworkUpdate(id: $id, name: $name, location: $location)"`
}

func (q UpdateRemoteNetwork) IsEmpty() bool {
	return q.Entity == nil
}

func (q UpdateRemoteNetwork) ToModel() *model.RemoteNetwork {
	if q.Entity == nil {
		return nil
	}

	return q.Entity.ToModel()
}
