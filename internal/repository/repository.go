package repository

import "fmt"

type RepositoryError error

type Repository interface {
	User() UserRepository
	UserSession() UserSessionRepository
}

var (
	RepositoryErrorNotFound          error = fmt.Errorf("record not found")
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
