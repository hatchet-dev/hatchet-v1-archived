package authorizer

import (
	"context"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/config/temporal"
	"github.com/hatchet-dev/hatchet/internal/temporal/server/authorizer/token"
	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/log"
)

type HatchetAuthorizer struct {
	config *temporal.Config

	tokenClaimMapper authorization.ClaimMapper
}

func NewHatchetAuthorizer(config *temporal.Config, authCfg *config.Authorization, logger log.Logger) *HatchetAuthorizer {
	keyProvider := token.NewHMACTokenKeyProvider(config.InternalAuthConfig)

	claimMapper := authorization.NewDefaultJWTClaimMapper(keyProvider, authCfg, logger)

	return &HatchetAuthorizer{config, claimMapper}
}

var decisionAllow = authorization.Result{Decision: authorization.DecisionAllow}
var decisionDeny = authorization.Result{Decision: authorization.DecisionDeny}

func (a *HatchetAuthorizer) Authorize(_ context.Context, claims *authorization.Claims,
	target *authorization.CallTarget) (authorization.Result, error) {

	if claims != nil && claims.Subject == "internal-admin" {
		return decisionAllow, nil
	}

	// Allow all requests
	return decisionDeny, nil
}

func (a *HatchetAuthorizer) GetClaims(authInfo *authorization.AuthInfo) (*authorization.Claims, error) {
	// wrap claims from default
	claims, err := a.tokenClaimMapper.GetClaims(authInfo)

	if err != nil {
		return nil, fmt.Errorf("error getting token claims: %s", err.Error())
	}

	// if TLS information is passed through, case on trusted CNs
	if authInfo.TLSSubject != nil {
		// TODO: pass this in via config
		switch authInfo.TLSSubject.CommonName {
		case "cluster":
			claims.Subject = "internal-admin"
			claims.System = authorization.RoleAdmin
		case "internal-admin":
			claims.Subject = "internal-admin"
			claims.System = authorization.RoleAdmin
		}
	}

	return claims, nil
}
