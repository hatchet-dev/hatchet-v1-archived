package authorizer

import (
	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
)

type myClaimMapper struct{}

func NewMyClaimMapper(_ *config.Config) authorization.ClaimMapper {
	return &myClaimMapper{}
}

func (c myClaimMapper) GetClaims(authInfo *authorization.AuthInfo) (*authorization.Claims, error) {
	// TODO: populate claims
	claims := authorization.Claims{}

	return &claims, nil
}
