package task

import (
	"log"
	"net/smtp"
)

type EmailTask struct {
	Auth         smtp.Auth
	Host         string
	Port         string
	From         string
	To           []string
	AlertMessage []byte
}

func NewEmailTask(identity, username, password, host, port, from string, to []string, alertMessage []byte) EmailTask {
	t := EmailTask{
		Host:         host,
		From:         from,
		To:           to,
		AlertMessage: alertMessage,
	}
	if port == "" {
		port = "25"
	}
	t.Port = port

	if username != "" && password != "" {
		auth := smtp.PlainAuth(identity, username, password, host)
		t.Auth = auth
	}

	return t
}

func (et EmailTask) Launch() {

	// msg := []byte("To: recipient@example.net\r\n" +
	// 	"Subject: discount Gophers!\r\n" +
	// 	"\r\n" +
	// 	"This is the email body.\r\n")
	err := smtp.SendMail(et.Host+":"+et.Port, et.Auth, et.From, et.To, et.AlertMessage)
	if err != nil {
		log.Fatal(err)
	}
}
