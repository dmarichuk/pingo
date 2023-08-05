package alert

import (
	"fmt"
	"log"
	"net/smtp"
	"pingo/job"
	"time"
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

func (a *EmailAlert) Send(msg []byte) {
	err := smtp.SendMail(a.Host+":"+a.Port, a.Auth, a.From, a.To, msg)
	if err != nil {
		log.Println(err)
	}
}

func (a *EmailAlert) GenerateMessage(j *job.Job) []byte {
	msg := fmt.Sprintf("From: %s\r\n", a.From)
	msg += fmt.Sprintf("Subject: %s\r\n", fmt.Sprintf("Job %s changed status to %s!", j.Name, j.Status))
	msg += fmt.Sprintf("\r\n%s\r\n", fmt.Sprintf("Job %s; Status: %s; Time: %s", j.Name, j.Status, j.TS.Format(time.RFC3339)))
	return []byte(msg)
}
