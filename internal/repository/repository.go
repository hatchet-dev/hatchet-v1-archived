package repository

import "fmt"

type RepositoryError error

type Repository interface {
	User() UserRepository
	UserSession() UserSessionRepository
	PasswordResetToken() PasswordResetTokenRepository
	VerifyEmailToken() VerifyEmailTokenRepository
	PersonalAccessToken() PersonalAccessTokenRepository
	Org() OrgRepository
	Team() TeamRepository
	GithubAppOAuth() GithubAppOAuthRepository
	GithubAppInstallation() GithubAppInstallationRepository
	GithubWebhook() GithubWebhookRepository
	GithubPullRequest() GithubPullRequestRepository
	Module() ModuleRepository
	ModuleValues() ModuleValuesRepository
	ModuleEnvVars() ModuleEnvVarsRepository
	ModuleRunQueue() ModuleRunQueueRepository
	ModuleMonitor() ModuleMonitorRepository
	Notification() NotificationRepository
	WorkerToken() WorkerTokenRepository
}

var (
	RepositoryErrorNotFound          error = fmt.Errorf("record not found")
	RepositoryNoRowsAffected         error = fmt.Errorf("no rows affected")
	RepositoryUniqueConstraintFailed error = fmt.Errorf("unique constraint failed")
)

func UnknownRepositoryError(err error) RepositoryError {
	return &RepositoryErrorUnknown{err}
}

type RepositoryErrorUnknown struct {
	err error
}

func (r *RepositoryErrorUnknown) Error() string {
	return fmt.Sprintf("unknown repository error: %v", r.err)
}

func BoolPointer(b bool) *bool {
	return &b
}
