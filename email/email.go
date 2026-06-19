package email

import (
	"fmt"
	"os"

	resend "github.com/resend/resend-go/v2"
)

var client *resend.Client
var fromEmail string

func init() {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		fmt.Println("WARN: RESEND_API_KEY not set — emails will be logged but not sent")
		return
	}
	client = resend.NewClient(apiKey)

	fromEmail = os.Getenv("RESEND_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev"
	}
}

func SendContactEmail(name, senderEmail, message string) error {
	if client == nil {
		fmt.Printf("EMAIL [not sent — no API key]: name=%s email=%s message=%s\n", name, senderEmail, message)
		return nil
	}

	params := &resend.SendEmailRequest{
		From:    fromEmail,
		To:      []string{"kylereynoldsdev@gmail.com"},
		ReplyTo: senderEmail,
		Subject: fmt.Sprintf("New contact from %s via your portfolio", name),
		Html: fmt.Sprintf(
			`<p><strong>Name:</strong> %s</p>
			 <p><strong>Email:</strong> %s</p>
			 <p><strong>Message:</strong></p>
			 <p>%s</p>`,
			name, senderEmail, message,
		),
	}

	_, err := client.Emails.Send(params)
	return err
}
