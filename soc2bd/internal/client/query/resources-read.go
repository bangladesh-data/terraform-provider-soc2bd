package query

import (
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/utils"
)

const CursorResources = "resourcesEndCursor"

type ReadResources struct {
	Resources `graphql:"resources(after: $resourcesEndCursor, first: $pageLimit)"`
}

func (r ReadResources) IsEmpty() bool {
	return len(r.Edges) == 0
}

type Resources struct {
	PaginatedResource[*ResourceEdge]
}

type ResourceEdge struct {
	Node *ResourceNode
}

func (r Resources) ToModel() []*model.Resource {
	return utils.Map[*ResourceEdge, *model.Resource](r.Edges, func(edge *ResourceEdge) *model.Resource {
		return edge.Node.ToModel()
	})
}
