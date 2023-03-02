package gorm

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NotificationRepository uses gorm.DB for querying the database
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository returns a DefaultNotificationRepository which uses
// gorm.DB for querying the database
func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &NotificationRepository{db}
}

func (repo *NotificationRepository) CreateNotificationInbox(inbox *models.NotificationInbox) (*models.NotificationInbox, repository.RepositoryError) {
	if err := repo.db.Omit(clause.Associations).Create(inbox).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return inbox, nil
}

func (repo *NotificationRepository) ReadNotificationInboxByTeamID(teamID string) (*models.NotificationInbox, repository.RepositoryError) {
	inbox := &models.NotificationInbox{}

	if err := repo.db.Where("team_id = ?", teamID).First(&inbox).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return inbox, nil
}

func (repo *NotificationRepository) UpdateNotificationInbox(inbox *models.NotificationInbox) (*models.NotificationInbox, repository.RepositoryError) {
	if err := repo.db.Omit(clause.Associations).Save(inbox).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return inbox, nil
}

func (repo *NotificationRepository) CreateNotification(notif *models.Notification) (*models.Notification, repository.RepositoryError) {
	if err := repo.db.Omit(clause.Associations).Create(notif).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return notif, nil
}

func (repo *NotificationRepository) ReadNotificationByID(teamID, id string) (*models.Notification, repository.RepositoryError) {
	notif := &models.Notification{}

	if err := repo.db.Preload("Module").Preload("MonitorResults").Preload("Runs").Where("team_id = ? AND id = ?", teamID, id).First(&notif).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return notif, nil
}

func (repo *NotificationRepository) ReadNotificationByNotificationID(teamID, notificationID string, opts *repository.ReadNotificationOpts) (*models.Notification, repository.RepositoryError) {
	notif := &models.Notification{}

	query := repo.db.Preload("Module").Where("team_id = ? AND notification_id = ?", teamID, notificationID)

	if opts.AutoResolved != nil {
		query = query.Where("auto_resolved = ?", *opts.AutoResolved)
	}

	if err := query.First(&notif).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return notif, nil
}

func (repo *NotificationRepository) ListNotifications(
	filterOpts *repository.ListNotificationOpts,
	opts ...repository.QueryOption,
) ([]*models.Notification, *repository.PaginatedResult, repository.RepositoryError) {
	var results []*models.Notification

	query := repo.db.Preload("Module").Model(&models.Notification{})

	if filterOpts.AutoResolved != nil {
		query = query.Where("auto_resolved = ?", *filterOpts.AutoResolved)
	}

	if filterOpts.Resolved != nil {
		query = query.Where("resolved = ?", *filterOpts.Resolved)
	}

	paginatedResult := &repository.PaginatedResult{}

	query = query.Scopes(queryutils.Paginate(opts, query, paginatedResult))

	if err := query.Find(&results).Error; err != nil {
		return nil, nil, err
	}

	return results, paginatedResult, nil
}

func (repo *NotificationRepository) ListNotificationsByTeamIDs(
	teamIDs []string,
	filterOpts *repository.ListNotificationOpts,
	opts ...repository.QueryOption) ([]*models.Notification, *repository.PaginatedResult, repository.RepositoryError) {
	if teamIDs == nil || len(teamIDs) == 0 {
		return nil, nil, repository.UnknownRepositoryError(fmt.Errorf("must pass in at least one team id to ListNotificationsByTeamIDs"))
	}

	var results []*models.Notification

	query := repo.db.Preload("Module").Model(&models.Notification{}).Where("team_id IN (?)", teamIDs)

	if filterOpts.AutoResolved != nil {
		query = query.Where("auto_resolved = ?", *filterOpts.AutoResolved)
	}

	if filterOpts.Resolved != nil {
		query = query.Where("resolved = ?", *filterOpts.Resolved)
	}

	paginatedResult := &repository.PaginatedResult{}

	query = query.Scopes(queryutils.Paginate(opts, query, paginatedResult))

	if err := query.Find(&results).Error; err != nil {
		return nil, nil, err
	}

	return results, paginatedResult, nil
}

func (repo *NotificationRepository) UpdateNotification(notif *models.Notification) (*models.Notification, repository.RepositoryError) {
	if err := repo.db.Omit(clause.Associations).Save(notif).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return notif, nil
}

func (repo *NotificationRepository) AppendModuleRun(notif *models.Notification, run *models.ModuleRun) (*models.Notification, repository.RepositoryError) {
	if err := repo.db.Model(notif).Omit("MonitorResults", "Runs.*", "Module").Association("Runs").Append(run); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return notif, nil
}

func (repo *NotificationRepository) AppendModuleRunMonitorResult(notif *models.Notification, result *models.ModuleMonitorResult) (*models.Notification, repository.RepositoryError) {
	if err := repo.db.Model(notif).Omit("MonitorResults.*", "Runs", "Module").Association("MonitorResults").Append(result); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return notif, nil
}

func (repo *NotificationRepository) DeleteNotification(notif *models.Notification) (*models.Notification, repository.RepositoryError) {
	del := repo.db.Delete(&notif)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return notif, nil
}
