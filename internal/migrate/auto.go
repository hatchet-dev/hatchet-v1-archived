package migrate

import (
	"github.com/hatchet-dev/hatchet/internal/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB, debug bool) error {
	instanceDB := db

	if debug {
		instanceDB = instanceDB.Debug()
	}

	return instanceDB.AutoMigrate(
		&models.User{},
		&models.UserSession{},
		&models.PersonalAccessToken{},
		&models.PasswordResetToken{},
		&models.VerifyEmailToken{},
		&models.Organization{},
		&models.OrganizationMember{},
		&models.OrganizationPolicy{},
		&models.OrganizationInviteLink{},
		&models.Team{},
		&models.TeamMember{},
		&models.TeamPolicy{},
		&models.GithubAppOAuth{},
		&models.GithubAppInstallation{},
		&models.GithubWebhook{},
		&models.GithubPullRequest{},
		&models.GithubPullRequestComment{},
		&models.Module{},
		&models.ModuleDeploymentConfig{},
		&models.ModuleRun{},
		&models.ModuleRunToken{},
		&models.ModuleRunConfig{},
		&models.ModuleValuesVersion{},
		&models.ModuleValues{},
	)
}
