package alert

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailAlert struct {
	Auth smtp.Auth
	Host string
	Port string
	From string
	To   []string
}

func NewEmailAlert(identity, username, password, host, port, from string, to []string) *EmailAlert {
	t := EmailAlert{
		Host: host,
		From: from,
		To:   to,
	}
	if port == "" {
		port = "25"
	}
	t.Port = port

	if username != "" && password != "" {
		auth := smtp.PlainAuth(identity, username, password, host)
		t.Auth = auth
	}

	return &t
}

func (a *EmailAlert) Send(msg string) {
	email := fmt.Sprintf("From: %s\r\n", a.From)
	email += fmt.Sprintf("Subject: Pingo Alert!\r\n")
	email += fmt.Sprintf("\r\n%s\r\n", msg)
	err := smtp.SendMail(a.Host+":"+a.Port, a.Auth, a.From, a.To, []byte(msg))
	if err != nil {
		log.Println(err)
	}
}
