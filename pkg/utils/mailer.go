package utils

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

func SendEmail(smtpCfg SMTPConfig, to []string, subject, body string) error {
	if smtpCfg.Host == "" {
		return fmt.Errorf("SMTP host is not configured")
	}

	msg := []byte("To: " + to[0] + "\r\n" +
		"From: " + smtpCfg.From + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", smtpCfg.User, smtpCfg.Password, smtpCfg.Host)
	addr := smtpCfg.Host + ":" + smtpCfg.Port

	fromAddr, err := mail.ParseAddress(smtpCfg.From)
	fromEmail := smtpCfg.From
	if err == nil {
		fromEmail = fromAddr.Address
	}

	return smtp.SendMail(addr, auth, fromEmail, to, msg)
}
