package models

import (
	"fmt"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSecurityPolicyModel(t *testing.T) {
	cases := []struct {
		policy model.SecurityPolicy

		expected interface{}
	}{
		{
			policy: model.SecurityPolicy{},
			expected: map[string]interface{}{
				attr.ID:   "",
				attr.Name: "",
			},
		},
		{
			policy: model.SecurityPolicy{
				ID:   "id",
				Name: "name",
			},
			expected: map[string]interface{}{
				attr.ID:   "id",
				attr.Name: "name",
			},
		},
	}

	for n, c := range cases {
		t.Run(fmt.Sprintf("case_%d", n), func(t *testing.T) {
			assert.Equal(t, c.expected, c.policy.ToTerraform())
		})
	}
}
