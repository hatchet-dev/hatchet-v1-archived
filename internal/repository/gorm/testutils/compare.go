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

func AssertUserSessionsEqual(t *testing.T, userSess1, userSess2 *models.UserSession) {
	AssertBaseEqual(t, &userSess1.Base, &userSess2.Base)

	userSess1Cp := *userSess1
	userSess2Cp := *userSess2

	userSess1Cp.Base = userSess2Cp.Base

	// check expires on within 10-second delta
	assert.InDelta(t, userSess1Cp.ExpiresAt.Unix(), userSess1Cp.ExpiresAt.Unix(), 1)

	userSess1Cp.ExpiresAt = userSess2.ExpiresAt

	assert.Equal(t, userSess1Cp, userSess2Cp, "user sessions should be equal")
}

func AssertPATsEqual(t *testing.T, pat1, pat2 *models.PersonalAccessToken) {
	AssertBaseEqual(t, &pat1.Base, &pat2.Base)

	pat1Cp := *pat1
	pat2Cp := *pat2

	pat1Cp.Base = pat2Cp.Base

	// check expires on within 10-second delta
	assert.InDelta(t, pat1Cp.Expires.Unix(), pat2Cp.Expires.Unix(), 1)

	pat1Cp.Expires = pat2Cp.Expires

	assert.Equal(t, pat1Cp, pat2Cp, "pats should be equal")
}
