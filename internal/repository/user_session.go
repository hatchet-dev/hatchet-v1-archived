package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// UserSessionRepository represents the set of queries on the UserSession model
type UserSessionRepository interface {
	CreateUserSession(session *models.UserSession) (*models.UserSession, RepositoryError)
	UpdateUserSession(session *models.UserSession) (*models.UserSession, RepositoryError)
	DeleteUserSession(session *models.UserSession) (*models.UserSession, RepositoryError)
	ReadUserSessionByKey(sessionKey string) (*models.UserSession, RepositoryError)
}
