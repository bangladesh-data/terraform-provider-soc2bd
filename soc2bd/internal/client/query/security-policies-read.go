package query

import (
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/utils"
)

const CursorPolicies = "policiesEndCursor"

type ReadSecurityPolicies struct {
	SecurityPolicies `graphql:"securityPolicies(after: $policiesEndCursor, first: $pageLimit)"`
}

func (q ReadSecurityPolicies) IsEmpty() bool {
	return len(q.Edges) == 0
}

type SecurityPolicies struct {
	PaginatedResource[*SecurityPolicyEdge]
}

type SecurityPolicyEdge struct {
	Node *gqlSecurityPolicy
}

func (q ReadSecurityPolicies) ToModel() []*model.SecurityPolicy {
	return utils.Map[*SecurityPolicyEdge, *model.SecurityPolicy](q.SecurityPolicies.Edges,
		func(edge *SecurityPolicyEdge) *model.SecurityPolicy {
			return edge.Node.ToModel()
		})
}
