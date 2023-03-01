package notifier

import "github.com/hatchet-dev/hatchet/internal/models"

type SendPasswordResetEmailOpts struct {
	Email string
	URL   string
}

type SendVerificationEmailOpts struct {
	Email string
	URL   string
}

type SendInviteLinkEmailOpts struct {
	Email            string
	URL              string
	OrganizationName string
	InviterAddress   string
}

type UserNotifier interface {
	SendPasswordResetEmail(opts *SendPasswordResetEmailOpts) error
	SendVerificationEmail(opts *SendVerificationEmailOpts) error
	SendInviteLinkEmail(opts *SendInviteLinkEmailOpts) error
}

type SendIncidentNotificationOpts struct {
	Users        []*models.User
	URL          string
	ModuleName   string
	Title        string
	Message      string
	Notification *models.Notification
}

type IncidentNotifier interface {
	SendIncidentNotification(opts *SendIncidentNotificationOpts) error
}
