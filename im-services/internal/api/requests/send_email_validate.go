package requests

type SendEmailRequest struct {
	Email     string `json:"email" validate:"required,email" `
	EmailType int    `json:"email_type" validate:"gte=1,lte=2"`
}
