package testutils

import (
	"testing"

	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/stretchr/testify/assert"
)

func AssertBaseEqual(t *testing.T, base1, base2 *models.Base) {
	assert.Equal(t, base1.CreatedAt.Unix(), base2.CreatedAt.Unix(), "created_at should be equal")
	assert.Equal(t, base1.UpdatedAt.Unix(), base2.UpdatedAt.Unix(), "created_at should be equal")
	assert.Equal(t, base1.ID, base2.ID, "ids should be equal")
}

func AssertUsersEqual(t *testing.T, user1, user2 *models.User) {
	AssertBaseEqual(t, &user1.Base, &user2.Base)

	user1Cp := *user1
	user2Cp := *user2

	user1Cp.Base = user2Cp.Base

	assert.Equal(t, user1Cp, user2Cp, "users should be equal")
}
