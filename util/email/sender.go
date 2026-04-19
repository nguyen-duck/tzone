package email

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
)

func SendOTP(to string, otp string, purpose string) error {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	port := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	user := strings.TrimSpace(os.Getenv("SMTP_USER"))
	pass := strings.ReplaceAll(strings.TrimSpace(os.Getenv("SMTP_PASS")), " ", "")
	from := strings.TrimSpace(os.Getenv("SMTP_FROM"))

	if host == "" || port == "" || user == "" || pass == "" {
		return fmt.Errorf("email service is not configured")
	}

	if from == "" {
		from = user
	}

	envelopeFrom := from
	fromHeader := from
	if parsed, err := mail.ParseAddress(from); err == nil {
		envelopeFrom = parsed.Address
		fromHeader = parsed.String()
	}

	subject := "TZone verification code"
	body := fmt.Sprintf("Your verification code for %s is: %s\nThis code expires in 5 minutes.", purpose, otp)
	msg := []byte("To: " + to + "\r\n" +
		"From: " + fromHeader + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", user, pass, host)
	addr := host + ":" + port

	return smtp.SendMail(addr, auth, envelopeFrom, []string{to}, msg)
}
