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
}

type SharedOpts struct {
	APIKey                 string
	SenderEmail            string
	RestrictedEmailDomains []string
}

func NewUserNotifier(opts *UserNotifierOpts) notifier.UserNotifier {
	return &UserNotifier{opts}
}

func (s *UserNotifier) SendPasswordResetEmail(opts *notifier.SendPasswordResetEmailOpts) error {
	if allowed := s.isAllowed(opts.Email); !allowed {
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
	if allowed := s.isAllowed(opts.Email); !allowed {
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

func (s *UserNotifier) isAllowed(target string) bool {
	if len(s.opts.RestrictedEmailDomains) == 0 {
		return true
	}

	targetComponents := strings.Split(target, "@")
	targetDomain := targetComponents[1]

	for _, domain := range s.opts.RestrictedEmailDomains {
		if domain == targetDomain {
			return true
		}
	}

	return false
}
