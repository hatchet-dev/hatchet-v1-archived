package notifier

type SendPasswordResetEmailOpts struct {
	Email string
	URL   string
}

type SendVerificationEmailOpts struct {
	Email string
	URL   string
}

type UserNotifier interface {
	SendPasswordResetEmail(opts *SendPasswordResetEmailOpts) error
	SendVerificationEmail(opts *SendVerificationEmailOpts) error
}
