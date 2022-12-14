package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

// GenerateTokenFromPAT creates a new JWT based on a personal access token model.
// Note that the personal access token model must include a signing secret and token ID
// already, otherwise the generation will fail.
func GenerateTokenFromPAT(pat *models.PersonalAccessToken) (string, error) {
	if len(pat.SigningSecret) <= 15 {
		return "", fmt.Errorf("signing secret must be at least 16 bytes in length")
	}

	rawTok, err := getJWTForUser(pat.UserID, pat.Base.ID)

	if err != nil {
		return "", err
	}

	return rawTok.encodeToken(pat.SigningSecret)
}

// GetPATFromEncoded returns a personal access token model based on the raw token. This method
// performs parsing and validatino on the raw token string, so the returned PAT can be considered
// valid without additional checks.
func GetPATFromEncoded(tokenString string, repo repository.PersonalAccessTokenRepository) (*models.PersonalAccessToken, error) {
	var pat *models.PersonalAccessToken

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (res interface{}, err error) {
		var signingSecret []byte

		signingSecret, pat, err = getSigningSecretAndPATFromToken(token, repo)

		return signingSecret, err
	})

	if err != nil {
		return nil, fmt.Errorf("could not parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return pat, nil
}

// IsPATValid parses and validates the raw token. It requires access to the repository because
// each token gets its own signing secret.
func IsPATValid(tokenString string, repo repository.PersonalAccessTokenRepository) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (res interface{}, err error) {
		signingSecret, _, err := getSigningSecretAndPATFromToken(token, repo)

		return signingSecret, err
	})

	return token != nil && err == nil, err
}

func getSigningSecretAndPATFromToken(token *jwt.Token, repo repository.PersonalAccessTokenRepository) ([]byte, *models.PersonalAccessToken, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	// we read the PAT by both the user ID and token ID fields, and use that to retrieve the
	// signing secret
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {

		return nil, nil, fmt.Errorf("claims could not be parsed before validate: error casting claims to jwt.MapClaims")
	}

	tokenID := claims["token_id"].(string)
	userID := claims["user_id"].(string)

	pat, err := repo.ReadPersonalAccessToken(userID, tokenID)

	if err != nil {
		return nil, nil, fmt.Errorf("claims could not be parsed before validate: personal access token could not be read: %v", err)
	}

	return pat.SigningSecret, pat, nil
}

type jwtClaims struct {
	*jwt.RegisteredClaims
	UserID  string `json:"user_id"`
	TokenID string `json:"token_id"`
}

func (t *jwtClaims) encodeToken(tokenSecret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(tokenSecret)
}

func getJWTForUser(userID, tokenID string) (*jwtClaims, error) {
	if userID == "" || !uuidutils.IsValidUUID(userID) {
		return nil, fmt.Errorf("user id must be a valid UUID")
	}

	if tokenID == "" || !uuidutils.IsValidUUID(tokenID) {
		return nil, fmt.Errorf("token id must be a valid UUID")
	}

	return &jwtClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		UserID:  userID,
		TokenID: tokenID,
	}, nil
}
