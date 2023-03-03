package token

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/hatchet-dev/hatchet/internal/config/temporal"
)

type TemporalClaims struct {
	*jwt.RegisteredClaims

	Permissions []string `json:"permissions"`
}

func GenerateInternalToken(authConfig *temporal.InternalAuthConfig) (string, error) {
	claims := &TemporalClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			Audience: authConfig.InternalTokenOpts.Audience,
			Issuer:   authConfig.InternalTokenOpts.Issuer,
			Subject:  "internal-admin",
		},
		Permissions: []string{"system:admin"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.Header["kid"] = "internal"

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(authConfig.InternalSigningKey)
}
