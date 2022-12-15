package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

type UserSessionRepository struct {
	db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) repository.UserSessionRepository {
	return &UserSessionRepository{db}
}

func (s *UserSessionRepository) CreateUserSession(session *models.UserSession) (*models.UserSession, repository.RepositoryError) {
	if err := s.db.Create(session).Error; err != nil {
		return nil, toRepoError(s.db, err)
	}

	return session, nil
}

func (s *UserSessionRepository) UpdateUserSession(session *models.UserSession) (*models.UserSession, repository.RepositoryError) {
	if err := s.db.Model(session).Where("Key = ?", session.Key).Updates(session).Error; err != nil {
		return nil, toRepoError(s.db, err)
	}

	return session, nil
}

func (s *UserSessionRepository) DeleteUserSession(session *models.UserSession) (*models.UserSession, repository.RepositoryError) {
	if err := s.db.Where("Key = ?", session.Key).Unscoped().Delete(session).Error; err != nil {
		return nil, toRepoError(s.db, err)
	}

	return session, nil
}

func (s *UserSessionRepository) ReadUserSessionByKey(sessionKey string) (*models.UserSession, repository.RepositoryError) {
	session := &models.UserSession{}

	if err := s.db.Where("Key = ?", sessionKey).First(session).Error; err != nil {
		return nil, toRepoError(s.db, err)
	}

	return session, nil
}
