package pats_test

import (
	"testing"
	"time"

	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
	"github.com/stretchr/testify/assert"
)

func AssertPATsEqual(t *testing.T, exp *types.PersonalAccessToken, got *types.PersonalAccessToken) {
	// assert that the uuid on the pat field is valid
	assert.True(t, uuidutils.IsValidUUID(got.ID))

	// check that the expiration is within 120 seconds
	assert.InDelta(t, got.Expires.Unix(), time.Now().Add(24*30*time.Hour).Unix(), 120)

	expResp := exp

	expResp.APIResourceMeta = got.APIResourceMeta // we already checked that these were equal
	expResp.Expires = got.Expires                 // we already checked that these were equal

	// assert the rest of the PAT object is valid
	assert.Equal(t, expResp, got, "pats should be equal")
}
