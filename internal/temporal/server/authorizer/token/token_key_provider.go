package token

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hatchet-dev/hatchet/internal/config/temporal"
)

type HMACTokenKeyProvider struct {
	authConfig *temporal.InternalAuthConfig
}

func NewHMACTokenKeyProvider(authConfig *temporal.InternalAuthConfig) *HMACTokenKeyProvider {
	return &HMACTokenKeyProvider{authConfig}
}

func (h *HMACTokenKeyProvider) EcdsaKey(alg string, kid string) (*ecdsa.PublicKey, error) {
	return nil, fmt.Errorf("unsupported key type ecdsa for: %s", alg)
}

func (h *HMACTokenKeyProvider) RsaKey(alg string, kid string) (*rsa.PublicKey, error) {
	return nil, fmt.Errorf("unsupported key type rsa for: %s", alg)
}

func (h *HMACTokenKeyProvider) HmacKey(alg string, kid string) ([]byte, error) {
	if kid == "internal" {
		return h.authConfig.InternalSigningKey, nil
	}

	return nil, fmt.Errorf("unsupported kid: %s", kid)
}

func (h *HMACTokenKeyProvider) SupportedMethods() []string {
	return []string{jwt.SigningMethodHS256.Name}
}

func (h *HMACTokenKeyProvider) Close() {
	return
}
