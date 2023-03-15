package noop

import (
	"github.com/hatchet-dev/hatchet/internal/logger"
	"github.com/hatchet-dev/hatchet/internal/notifier"
)

type UserNotifier struct {
	l *logger.Logger
}

func NewNoOpUserNotifier(l *logger.Logger) notifier.UserNotifier {
	return &UserNotifier{l}
}

func (s *UserNotifier) GetID() string {
	return "noop"
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
	l *logger.Logger
}

func NewNoOpIncidentNotifier(l *logger.Logger) notifier.IncidentNotifier {
	return &IncidentNotifier{l}
}

func (s *IncidentNotifier) SendIncidentNotification(opts *notifier.SendIncidentNotificationOpts) error {
	s.l.Error().Msgf("New incident %s: %s. Visit %s for more information.", opts.Title, opts.Message, opts.URL)

	return nil
}
