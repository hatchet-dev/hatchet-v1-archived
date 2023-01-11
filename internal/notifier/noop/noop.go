package noop

import (
	"github.com/hatchet-dev/hatchet/internal/notifier"
)

type UserNotifier struct {
}

func NewNoOpUserNotifier() notifier.UserNotifier {
	return &UserNotifier{}
}

func (s *UserNotifier) SendPasswordResetEmail(opts *notifier.SendPasswordResetEmailOpts) error {
	return nil
}

func (s *UserNotifier) SendVerificationEmail(opts *notifier.SendVerificationEmailOpts) error {
	return nil
}
