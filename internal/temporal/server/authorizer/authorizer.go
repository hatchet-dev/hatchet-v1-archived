package authorizer

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hatchet-dev/hatchet/internal/config/temporal"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
	"github.com/hatchet-dev/hatchet/internal/temporal/server/authorizer/token"

	hatchettoken "github.com/hatchet-dev/hatchet/internal/auth/token"
	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/log"
)

var permittedWorkerAPIs map[string]bool = map[string]bool{
	"/temporal.api.workflowservice.v1.WorkflowService/GetSystemInfo": true,
}

type HatchetAuthorizer struct {
	config *temporal.Config

	keyProvider      authorization.TokenKeyProvider
	tokenClaimMapper authorization.ClaimMapper
	logger           log.Logger
}

func NewHatchetAuthorizer(config *temporal.Config, authCfg *config.Authorization, logger log.Logger) *HatchetAuthorizer {
	keyProvider := token.NewHMACTokenKeyProvider(config.InternalAuthConfig)

	claimMapper := authorization.NewDefaultJWTClaimMapper(keyProvider, authCfg, logger)

	return &HatchetAuthorizer{config, keyProvider, claimMapper, logger}
}

var decisionAllow = authorization.Result{Decision: authorization.DecisionAllow}
var decisionDeny = authorization.Result{Decision: authorization.DecisionDeny}

func (a *HatchetAuthorizer) Authorize(_ context.Context, claims *authorization.Claims,
	target *authorization.CallTarget) (authorization.Result, error) {

	if claims != nil && claims.Subject == "internal-admin" {
		return decisionAllow, nil
	}

	if claims != nil && claims.Subject != "" {
		// ensure the UUID is valid
		if !uuidutils.IsValidUUID(claims.Subject) {
			return authorization.Result{
				Decision: authorization.DecisionDeny,
				Reason:   fmt.Sprintf("uuid is not valid: %s", claims.Subject),
			}, nil
		}

		// ensure the team exists
		team, err := a.config.DB.Repository.Team().ReadTeamByID(claims.Subject)

		if err != nil {
			return authorization.Result{
				Decision: authorization.DecisionDeny,
				Reason:   fmt.Sprintf("could not read team passed as subject: %s", err),
			}, nil
		}

		if target.Namespace != "" {
			// ensure the target namespace matches the team namespace
			if target.Namespace == team.ID {
				return decisionAllow, nil
			}

			denyReason := fmt.Sprintf("target namespace %s does not match team id %s", target.Namespace, team.ID)

			a.logger.Error(denyReason)

			return authorization.Result{
				Decision: authorization.DecisionDeny,
				Reason:   denyReason,
			}, nil
		}

		// if there's no namespace, check against permitted calls
		if permitted, ok := permittedWorkerAPIs[target.APIName]; ok && permitted {
			return decisionAllow, nil
		}

		denyReason := fmt.Sprintf("api call %s is not in permitted api call list", target.APIName)

		a.logger.Error(denyReason)

		return authorization.Result{
			Decision: authorization.DecisionDeny,
			Reason:   denyReason,
		}, nil
	}

	// Deny any requests that don't match the internal admin user
	return decisionDeny, nil
}

func (a *HatchetAuthorizer) GetClaims(authInfo *authorization.AuthInfo) (*authorization.Claims, error) {
	claims := &authorization.Claims{}

	// wrap claims from default
	if authInfo.AuthToken != "" {
		parts := strings.Split(authInfo.AuthToken, " ")

		if len(parts) != 2 {
			return nil, fmt.Errorf("unexpected authorization token format")
		}

		if !strings.EqualFold(parts[0], authorizationBearer) {
			return nil, fmt.Errorf("unexpected name in authorization token")
		}

		rawToken := parts[1]

		kind, err := a.getTokenKind(rawToken)

		if err != nil {
			return nil, err
		}

		switch kind {
		case TemporalTokenKindInternal:
			claims, err = a.tokenClaimMapper.GetClaims(authInfo)

			if err != nil {
				return nil, fmt.Errorf("error getting token claims: %s", err.Error())
			}

		case TemporalTokenKindWorker:
			wtRepo := a.config.DB.Repository.WorkerToken()

			wt, err := hatchettoken.GetWTFromEncoded(rawToken, wtRepo, &a.config.InternalAuthConfig.InternalTokenOpts)

			if err != nil {
				return nil, fmt.Errorf("error getting token claims: %s", err.Error())
			}

			claims.Subject = wt.TeamID
			claims.Namespaces = map[string]authorization.Role{
				wt.TeamID: authorization.RoleAdmin,
			}
		}
	}

	// if TLS information is passed through, case on trusted CNs
	if authInfo.TLSSubject != nil {
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

type TemporalTokenKind string

const (
	TemporalTokenKindInternal TemporalTokenKind = "internal"
	TemporalTokenKindWorker   TemporalTokenKind = "worker"
)

const authorizationBearer = "bearer"

func (a *HatchetAuthorizer) getTokenKind(rawToken string) (TemporalTokenKind, error) {
	var kind TemporalTokenKind
	var rawKID string

	jwt.Parse(rawToken, func(token *jwt.Token) (res interface{}, err error) {
		var ok bool
		rawKID, ok = token.Header["kid"].(string)

		if !ok {
			return nil, fmt.Errorf("Not validated")
		}

		switch rawKID {
		case "internal":
			kind = TemporalTokenKindInternal
		case "worker":
			kind = TemporalTokenKindWorker
		}

		return nil, fmt.Errorf("Not validated")
	})

	if kind == "" {
		return "", fmt.Errorf("%s is not a supported token kind", rawKID)
	}

	return kind, nil
}
