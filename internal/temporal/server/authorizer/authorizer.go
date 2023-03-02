package authorizer

import (
	"context"

	"go.temporal.io/server/common/authorization"
)

type myAuthorizer struct{}

func NewMyAuthorizer() authorization.Authorizer {
	return &myAuthorizer{}
}

var decisionAllow = authorization.Result{Decision: authorization.DecisionAllow}
var decisionDeny = authorization.Result{Decision: authorization.DecisionDeny}

func (a *myAuthorizer) Authorize(_ context.Context, claims *authorization.Claims,
	target *authorization.CallTarget) (authorization.Result, error) {
	// Allow all requests
	return decisionAllow, nil
}
