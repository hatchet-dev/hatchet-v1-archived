package token

import (
	"time"

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
			Audience:  authConfig.InternalTokenOpts.Audience,
			Issuer:    authConfig.InternalTokenOpts.Issuer,
			Subject:   "internal-admin",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(90 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		Permissions: []string{"system:admin"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.Header["kid"] = "internal"

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(authConfig.InternalSigningKey)
}
