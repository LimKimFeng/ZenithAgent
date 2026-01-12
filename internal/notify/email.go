package notify

import (
	"fmt"
	"net/smtp"
)

type EmailNotifier struct {
	User      string
	Password  string
	Recipient string
	Host      string
	Port      string
}

func NewEmailNotifier(user, password, recipient string) *EmailNotifier {
	return &EmailNotifier{
		User:      user,
		Password:  password,
		Recipient: recipient,
		Host:      "smtp.gmail.com", // Default
		Port:      "587",           // Default
	}
}

func (n *EmailNotifier) SendDailyReport() error {
	body := "Daily Report Triggered. Checking stats.json..."
	return n.SendEmail("Daily Status Report", body)
}

func (n *EmailNotifier) SendEmail(subject, body string) error {
	if n.User == "" || n.Password == "" || n.Recipient == "" {
		return fmt.Errorf("email configuration missing")
	}

	addr := fmt.Sprintf("%s:%s", n.Host, n.Port)
	auth := smtp.PlainAuth("", n.User, n.Password, n.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: ZenithAgent: %s\r\n"+
		"\r\n"+
		"%s\r\n", n.Recipient, subject, body))

	return smtp.SendMail(addr, auth, n.User, []string{n.Recipient}, msg)
}
