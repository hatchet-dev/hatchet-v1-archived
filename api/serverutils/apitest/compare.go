package apitest

import (
	"testing"

	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
	"github.com/stretchr/testify/assert"
)

func AssertAPIResourcesEqual(t *testing.T, exp *types.APIResourceMeta, got *types.APIResourceMeta) {
	assert.True(t, uuidutils.IsValidUUID(exp.ID))
	assert.True(t, uuidutils.IsValidUUID(got.ID))

	assert.Equal(t, exp.ID, got.ID, "ids should be equal")
}
