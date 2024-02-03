package utils

import (
	"example/web-service-gin/internal/database"
	"fmt"
	"net/smtp"
)

func SendEmail(to, subject, body string) bool {
	// // Sender data.
	// from := "kales.ali1963@gmail.com"
	// password := "tsnmouasolzintlb"

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", database.Email, database.EmailPassword, smtpHost)

	message := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, database.Email, []string{to}, message)
	if err != nil {
		return false
	}

	return true
}
