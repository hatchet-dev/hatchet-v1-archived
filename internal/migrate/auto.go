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
		&models.Organization{},
		&models.OrganizationMember{},
		&models.OrganizationPolicy{},
	)
}
