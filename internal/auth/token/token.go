package token

import (
	"fmt"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type TokenOpts struct {
	Issuer   string
	Audience []string
}

// GenerateTokenFromPAT creates a new JWT based on a personal access token model.
// Note that the personal access token model must include a signing secret and token ID
// already, otherwise the generation will fail.
func GenerateTokenFromPAT(pat *models.PersonalAccessToken, opts *TokenOpts) (string, error) {
	// we enforce that issuer and audience are valid URLs (with schema)
	if _, err := url.Parse(opts.Issuer); err != nil {
		return "", fmt.Errorf("could not parse issuer: %v", err)
	}

	for _, aud := range opts.Audience {
		if _, err := url.Parse(aud); err != nil {
			return "", fmt.Errorf("could not parse aud %s: %v", aud, err)
		}
	}

	if len(pat.SigningSecret) <= 15 {
		return "", fmt.Errorf("signing secret must be at least 16 bytes in length")
	}

	rawTok, err := getJWTForUser(pat.UserID, pat.Base.ID, opts)

	if err != nil {
		return "", err
	}

	return rawTok.encodeToken(pat.SigningSecret)
}

// GetPATFromEncoded returns a personal access token model based on the raw token. This method
// performs parsing and validatino on the raw token string, so the returned PAT can be considered
// valid without additional checks.
func GetPATFromEncoded(tokenString string, repo repository.PersonalAccessTokenRepository, opts *TokenOpts) (*models.PersonalAccessToken, error) {
	var pat *models.PersonalAccessToken

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (res interface{}, err error) {
		var signingSecret []byte

		signingSecret, pat, err = getSigningSecretAndPATFromToken(token, repo, opts)

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
func IsPATValid(tokenString string, repo repository.PersonalAccessTokenRepository, opts *TokenOpts) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (res interface{}, err error) {
		signingSecret, _, err := getSigningSecretAndPATFromToken(token, repo, opts)

		return signingSecret, err
	})

	return token != nil && err == nil, err
}

func getSigningSecretAndPATFromToken(token *jwt.Token, repo repository.PersonalAccessTokenRepository, opts *TokenOpts) ([]byte, *models.PersonalAccessToken, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	// we read the PAT by both the user ID and token ID fields, and use that to retrieve the
	// signing secret
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, nil, fmt.Errorf("claims could not be parsed before validate: error casting claims to jwt.MapClaims")
	}

	// check that the token kind is a PAT token
	if claims["kind"].(string) != string(jwtClaimKindPAT) {
		return nil, nil, fmt.Errorf("claim kind was not pat")
	}

	if claims["iss"].(string) != opts.Issuer {
		return nil, nil, fmt.Errorf("issuer was not %s", opts.Issuer)
	}

	audience := claims["aud"].([]interface{})
	matchedAud := false

	for _, aud1 := range audience {
		for _, aud2 := range opts.Audience {
			if matchedAud = aud1.(string) == aud2; matchedAud {
				break
			}
		}
	}

	if !matchedAud {
		return nil, nil, fmt.Errorf("did not find an audience match for the token")
	}

	tokenID := claims["token_id"].(string)
	userID := claims["sub"].(string)

	pat, err := repo.ReadPersonalAccessToken(userID, tokenID)

	if err != nil {
		return nil, nil, fmt.Errorf("claims could not be parsed before validate: personal access token could not be read: %v", err)
	}

	return pat.SigningSecret, pat, nil
}

type JWTClaimKind string

const (
	jwtClaimKindPAT JWTClaimKind = "pat"
	jwtClaimKindAPI JWTClaimKind = "api"
)

// exported for testing purposes. This should NOT be used by callers.
type JWTClaims struct {
	*jwt.RegisteredClaims

	// Kind represents the type of JWT token for this server
	Kind    JWTClaimKind `json:"kind"`
	TokenID string       `json:"token_id"`
}

func (t *JWTClaims) encodeToken(tokenSecret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(tokenSecret)
}

func getJWTForUser(userID, tokenID string, opts *TokenOpts) (*JWTClaims, error) {
	if userID == "" || !uuidutils.IsValidUUID(userID) {
		return nil, fmt.Errorf("user id must be a valid UUID")
	}

	if tokenID == "" || !uuidutils.IsValidUUID(tokenID) {
		return nil, fmt.Errorf("token id must be a valid UUID")
	}

	return &JWTClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   userID,
			Audience:  opts.Audience,
			Issuer:    opts.Issuer,
		},
		TokenID: tokenID,
		Kind:    jwtClaimKindPAT,
	}, nil
}
