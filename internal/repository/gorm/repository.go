package gorm

import (
	"errors"

	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type GormRepository struct {
	user        repository.UserRepository
	userSession repository.UserSessionRepository
	pat         repository.PersonalAccessTokenRepository
	prt         repository.PasswordResetTokenRepository
	vet         repository.VerifyEmailTokenRepository
	org         repository.OrgRepository
	team        repository.TeamRepository
}

func (t *GormRepository) User() repository.UserRepository {
	return t.user
}

func (t *GormRepository) UserSession() repository.UserSessionRepository {
	return t.userSession
}

func (t *GormRepository) PersonalAccessToken() repository.PersonalAccessTokenRepository {
	return t.pat
}

func (t *GormRepository) PasswordResetToken() repository.PasswordResetTokenRepository {
	return t.prt
}

func (t *GormRepository) VerifyEmailToken() repository.VerifyEmailTokenRepository {
	return t.vet
}

func (t *GormRepository) Org() repository.OrgRepository {
	return t.org
}

func (t *GormRepository) Team() repository.TeamRepository {
	return t.team
}

// NewRepository returns a Repository which persists users in memory
// and accepts a parameter that can trigger read/write errors
func NewRepository(db *gorm.DB, key *[32]byte) repository.Repository {
	return &GormRepository{
		user:        NewUserRepository(db),
		userSession: NewUserSessionRepository(db),
		pat:         NewPersonalAccessTokenRepository(db, key),
		prt:         NewPasswordResetTokenRepository(db),
		vet:         NewVerifyEmailTokenRepository(db),
		org:         NewOrgRepository(db),
		team:        NewTeamRepository(db),
	}
}

func toRepoError(db *gorm.DB, err error) repository.RepositoryError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.RepositoryErrorNotFound
	}

	switch db.Dialector.Name() {
	case "sqlite":
		if err, ok := err.(sqlite3.Error); ok {
			if err.ExtendedCode == sqlite3.ErrConstraintUnique {
				return repository.RepositoryUniqueConstraintFailed

			}
		}
	case "postgres":
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == pgerrcode.UniqueViolation {
				return repository.RepositoryUniqueConstraintFailed
			}
		}
	}

	return repository.UnknownRepositoryError(err)
}
