package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type Subject string

const (
	User Subject = "user"
	API  Subject = "api"
)

type TokenGeneratorConf struct {
	TokenSecret string
}

type JWT struct {
	*jwt.RegisteredClaims
	UserID  string `json:"user_id"`
	TokenID string `json:"token_id"`
}

func GetJWTForUser(userID, tokenID string) (*JWT, error) {
	if userID == "" || !uuidutils.IsValidUUID(userID) {
		return nil, fmt.Errorf("user id must be a valid UUID")
	}

	if tokenID == "" || !uuidutils.IsValidUUID(tokenID) {
		return nil, fmt.Errorf("token id must be a valid UUID")
	}

	return &JWT{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		UserID:  userID,
		TokenID: tokenID,
	}, nil
}

func (t *JWT) EncodeToken(conf *TokenGeneratorConf) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(conf.TokenSecret))
}

func GetPATFromEncoded(tokenString string, repo repository.PersonalAccessTokenRepository) (*models.PersonalAccessToken, error) {
	var pat *models.PersonalAccessToken

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (res interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// we read the PAT by both the user ID and token ID fields, and use that to retrieve the
		// signing secret
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {

			return nil, fmt.Errorf("claims could not be parsed before validate: error casting claims to jwt.MapClaims")
		}

		tokenID := claims["token_id"].(string)
		userID := claims["user_id"].(string)

		pat, err = repo.ReadPersonalAccessToken(userID, tokenID)

		if err != nil {
			return nil, fmt.Errorf("claims could not be parsed before validate: personal access token could not be read: %v", err)
		}

		return pat.SigningSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return pat, nil
}

// TODO (abelanger5): remove all this logic for encrypting the secret key from the token package
// into Encrypt() method in the models package. Models are responsible for their own encryption.
func GenerateTokenForPAT(pat *models.PersonalAccessToken) (string, error) {
	rawTok, err := GetJWTForUser(pat.UserID, pat.Base.ID)

	if err != nil {
		return "", err
	}

	return rawTok.EncodeToken(&TokenGeneratorConf{
		TokenSecret: string(pat.SigningSecret),
	})
}
