package client

import (
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestUnwrapError(t *testing.T) {
	err := client.NewAPIError(errBadRequest, "read", "resource")

	assert.Equal(t, errBadRequest, err.Unwrap())
}
