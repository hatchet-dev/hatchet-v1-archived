package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// UserSessionRepository represents the set of queries on the UserSession model
type UserSessionRepository interface {
	CreateUserSession(session *models.UserSession) (*models.UserSession, error)
	UpdateUserSession(session *models.UserSession) (*models.UserSession, error)
	DeleteUserSession(session *models.UserSession) (*models.UserSession, error)
	ReadUserSessionByKey(sessionKey string) (*models.UserSession, error)
}
