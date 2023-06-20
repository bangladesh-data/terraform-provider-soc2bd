package query

import (
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/utils"
)

const CursorUsers = "usersEndCursor"

type ReadUsers struct {
	Users `graphql:"users(after: $usersEndCursor, first: $pageLimit)"`
}

func (q ReadUsers) IsEmpty() bool {
	return len(q.Edges) == 0
}

type Users struct {
	PaginatedResource[*UserEdge]
}

type UserEdge struct {
	Node *gqlUser
}

func (u Users) ToModel() []*model.User {
	return utils.Map[*UserEdge, *model.User](u.Edges, func(edge *UserEdge) *model.User {
		return edge.Node.ToModel()
	})
}
