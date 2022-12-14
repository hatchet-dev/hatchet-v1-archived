package token_test

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetJWTForUser(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		rawTok, err := token.GenerateTokenFromPAT(testutils.PATModels[0])

		assert.Nil(t, err, "error should be nil")

		isValid, err := token.IsPATValid(rawTok, conf.Repository.PersonalAccessToken())

		assert.Nil(t, err, "error should be nil")
		assert.True(t, isValid, "token should be valid")

		return nil
	}, testutils.InitUsers, testutils.InitPATs)
}

type jwtClaims struct {
	*jwt.RegisteredClaims
	UserID  string `json:"user_id"`
	TokenID string `json:"token_id"`
}

func TestIsJWTValid(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		originalTok, err := token.GenerateTokenFromPAT(testutils.PATModels[0])

		assert.Nil(t, err, "error should be nil")

		isValid, err := token.IsPATValid(originalTok, conf.Repository.PersonalAccessToken())

		assert.Nil(t, err, "error should be nil")
		assert.True(t, isValid, "token should be valid")

		// modify a single claim string and ensure that the token is no longer valid
		modifiedTokParts := strings.Split(originalTok, ".")

		decodedBytes, err := base64.RawStdEncoding.DecodeString(modifiedTokParts[1])

		assert.Nil(t, err, "error should be nil")

		claims := new(jwtClaims)

		err = json.Unmarshal(decodedBytes, &claims)
		assert.Nil(t, err, "error should be nil")

		claims.RegisteredClaims.Audience = jwt.ClaimStrings{"h"}

		newClaimsStr, err := json.Marshal(claims)
		assert.Nil(t, err, "error should be nil")

		newEncodedBytes := base64.RawStdEncoding.EncodeToString(newClaimsStr)

		modifiedTokParts[1] = newEncodedBytes

		newInvalidTok := strings.Join(modifiedTokParts, ".")

		isValid, err = token.IsPATValid(newInvalidTok, conf.Repository.PersonalAccessToken())

		assert.NotNil(t, err, "error should not be nil")
		assert.False(t, isValid, "invalid token should not be valid")
		assert.ErrorContains(t, err, "signature is invalid")

		// modify the header with a new algorithm and ensure that the token is no longer valid
		newAlgHeader := "{\"alg\": \"none\",\"typ\":\"JWT\"}"

		newEncodedBytes = base64.RawStdEncoding.EncodeToString([]byte(newAlgHeader))

		modifiedTokParts2 := strings.Split(originalTok, ".")

		modifiedTokParts2[0] = newEncodedBytes

		newInvalidTok = strings.Join(modifiedTokParts2, ".")
		isValid, err = token.IsPATValid(newInvalidTok, conf.Repository.PersonalAccessToken())

		assert.NotNil(t, err, "error should not be nil")
		assert.False(t, isValid, "invalid token should not be valid")
		assert.ErrorContains(t, err, "Unexpected signing method: none")

		return nil
	}, testutils.InitUsers, testutils.InitPATs)
}

func TestGenerateTokenFromPATInvalidSigningSecret(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		modelCp := testutils.PATModels[0]

		modelCp.SigningSecret = []byte("abcd")

		_, err := token.GenerateTokenFromPAT(modelCp)

		assert.NotNil(t, err, "error should not be nil")
		assert.ErrorContains(t, err, "signing secret must be at least 16 bytes in length")

		return nil
	}, testutils.InitUsers, testutils.InitPATs)
}

func TestGetPATFromEncoded(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		rawTok, err := token.GenerateTokenFromPAT(testutils.PATModels[0])

		assert.Nil(t, err, "error should be nil")

		// recover the PAT model
		pat, err := token.GetPATFromEncoded(rawTok, conf.Repository.PersonalAccessToken())

		assert.Nil(t, err, "error should be nil")

		expRes := testutils.PATModels[0]

		// reset the time fields to check equality
		expRes.Base.CreatedAt = pat.CreatedAt
		expRes.Base.UpdatedAt = pat.UpdatedAt
		expRes.Expires = pat.Expires

		assert.Equal(t, expRes, pat, "pat models should be equal")

		return nil
	}, testutils.InitUsers, testutils.InitPATs)

}
