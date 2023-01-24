package models

import (
	"github.com/hatchet-dev/hatchet/api/v1/types"
)

type GithubAppInstallation struct {
	Base

	GithubAppOAuthID string
	GithubAppOAuth   GithubAppOAuth `gorm:"foreignKey:GithubAppOAuthID"`

	AccountName             string
	AccountAvatarURL        string
	AccountID               int64
	InstallationID          int64
	InstallationSettingsURL string
}

func (g *GithubAppInstallation) ToAPIType() *types.GithubAppInstallation {
	return &types.GithubAppInstallation{
		APIResourceMeta:         g.ToAPITypeMetadata(),
		InstallationID:          g.InstallationID,
		InstallationSettingsURL: g.InstallationSettingsURL,
		AccountName:             g.AccountName,
		AccountAvatarURL:        g.AccountAvatarURL,
	}
}
