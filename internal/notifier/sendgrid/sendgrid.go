package sendgrid

import (
	"fmt"
	"strings"

	"github.com/hatchet-dev/hatchet/internal/notifier"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type UserNotifier struct {
	opts *UserNotifierOpts
}

type UserNotifierOpts struct {
	*SharedOpts
	PWResetTemplateID     string
	VerifyEmailTemplateID string
	InviteLinkTemplateID  string
}

type SharedOpts struct {
	APIKey                 string
	SenderEmail            string
	RestrictedEmailDomains []string
}

func NewUserNotifier(opts *UserNotifierOpts) notifier.UserNotifier {
	return &UserNotifier{opts}
}

func (s *UserNotifier) GetID() string {
	return "sendgrid"
}

func (s *UserNotifier) SendPasswordResetEmail(opts *notifier.SendPasswordResetEmailOpts) error {
	if allowed := s.opts.isAllowed(opts.Email); !allowed {
		return fmt.Errorf("target email %s is not in restricted domain group", opts.Email)
	}

	request := sendgrid.GetRequest(s.opts.APIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"

	sgMail := &mail.SGMailV3{
		Personalizations: []*mail.Personalization{
			{
				To: []*mail.Email{
					{
						Address: opts.Email,
					},
				},
				DynamicTemplateData: map[string]interface{}{
					"url":   opts.URL,
					"email": opts.Email,
				},
			},
		},
		From: &mail.Email{
			Address: s.opts.SenderEmail,
			Name:    "Hatchet",
		},
		TemplateID: s.opts.PWResetTemplateID,
	}

	request.Body = mail.GetRequestBody(sgMail)

	_, err := sendgrid.API(request)

	return err
}

func (s *UserNotifier) SendVerificationEmail(opts *notifier.SendVerificationEmailOpts) error {
	if allowed := s.opts.isAllowed(opts.Email); !allowed {
		return fmt.Errorf("target email %s is not in restricted domain group", opts.Email)
	}

	request := sendgrid.GetRequest(s.opts.APIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"

	sgMail := &mail.SGMailV3{
		Personalizations: []*mail.Personalization{
			{
				To: []*mail.Email{
					{
						Address: opts.Email,
					},
				},
				DynamicTemplateData: map[string]interface{}{
					"url":   opts.URL,
					"email": opts.Email,
				},
			},
		},
		From: &mail.Email{
			Address: s.opts.SenderEmail,
			Name:    "Hatchet",
		},
		TemplateID: s.opts.VerifyEmailTemplateID,
	}

	request.Body = mail.GetRequestBody(sgMail)

	_, err := sendgrid.API(request)

	return err
}

func (s *UserNotifier) SendInviteLinkEmail(opts *notifier.SendInviteLinkEmailOpts) error {
	if allowed := s.opts.isAllowed(opts.Email); !allowed {
		return fmt.Errorf("target email %s is not in restricted domain group", opts.Email)
	}

	request := sendgrid.GetRequest(s.opts.APIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"

	sgMail := &mail.SGMailV3{
		Personalizations: []*mail.Personalization{
			{
				To: []*mail.Email{
					{
						Address: opts.Email,
					},
				},
				DynamicTemplateData: map[string]interface{}{
					"url":               opts.URL,
					"email":             opts.Email,
					"organization_name": opts.OrganizationName,
					"inviter_address":   opts.InviterAddress,
				},
			},
		},
		From: &mail.Email{
			Address: s.opts.SenderEmail,
			Name:    "Hatchet",
		},
		TemplateID: s.opts.InviteLinkTemplateID,
	}

	request.Body = mail.GetRequestBody(sgMail)

	_, err := sendgrid.API(request)

	return err
}

type IncidentNotifier struct {
	opts *IncidentNotifierOpts
}

type IncidentNotifierOpts struct {
	*SharedOpts
	IncidentTemplateID string
}

func NewIncidentNotifier(opts *IncidentNotifierOpts) notifier.IncidentNotifier {
	return &IncidentNotifier{opts}
}

func (s *IncidentNotifier) SendIncidentNotification(opts *notifier.SendIncidentNotificationOpts) error {
	for _, user := range opts.Users {
		if allowed := s.opts.isAllowed(user.Email); !allowed {
			return fmt.Errorf("target email %s is not in restricted domain group", user.Email)
		}
	}

	request := sendgrid.GetRequest(s.opts.APIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"

	personalizations := make([]*mail.Personalization, 0)

	templData := map[string]interface{}{
		"title":      opts.Title,
		"message":    opts.Message,
		"subject":    fmt.Sprintf("[%s] %s", opts.ModuleName, opts.Title),
		"created_at": fmt.Sprintf("%s", opts.Notification.CreatedAt.Format("Jan 2, 2006 at 3:04pm (MST)")),
		"url":        opts.URL,
	}

	for _, user := range opts.Users {
		personalizations = append(personalizations, &mail.Personalization{
			To: []*mail.Email{
				{
					Address: user.Email,
				},
			},
			DynamicTemplateData: templData,
		})
	}

	sgMail := &mail.SGMailV3{
		Personalizations: personalizations,
		From: &mail.Email{
			Address: s.opts.SenderEmail,
			Name:    "Porter Notifications",
		},
		TemplateID: s.opts.IncidentTemplateID,
	}

	request.Body = mail.GetRequestBody(sgMail)

	_, err := sendgrid.API(request)

	return err

}

func (s *SharedOpts) isAllowed(target string) bool {
	if len(s.RestrictedEmailDomains) == 0 {
		return true
	}

	targetComponents := strings.Split(target, "@")
	targetDomain := targetComponents[1]

	for _, domain := range s.RestrictedEmailDomains {
		if domain == targetDomain {
			return true
		}
	}

	return false
}
