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

func (s *UserNotifier) SendInviteLinkEmail(opts *notifier.SendInviteLinkEmailOpts) error {
	return nil
}

type IncidentNotifier struct {
}

func NewNoOpIncidentNotifier() notifier.IncidentNotifier {
	return &IncidentNotifier{}
}

func (s *IncidentNotifier) SendIncidentNotification(opts *notifier.SendIncidentNotificationOpts) error {
	return nil
}
