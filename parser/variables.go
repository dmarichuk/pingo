package parser

import (
	"log"
	"os"
	"reflect"
	"strings"
)

type Variables struct {
	// Telegram
	TelegramBotToken string `yaml:"telegram_bot_token,omitempty"`
	TelegramChatID   string `yaml:"telegram_chat_id,omitempty"`
	// SMTP
	SmtpIdentity string   `yaml:"smtp_identity,omitempty"`
	SmtpHost     string   `yaml:"smtp_host,omitempty"`
	SmtpPort     string   `yaml:"smtp_port,omitempty"`
	SmtpUsername string   `yaml:"smtp_username,omitempty"`
	SmtpPassword string   `yaml:"smtp_password,omitempty"`
	SmtpFrom     string   `yaml:"smtp_from,omitempty"`
	SmtpTo       []string `yaml:"smtp_recipients,omitempty"`
}

func (vars *Variables) PostParse() {
	// Check if there are any ${{...}} enviroment variables
	v := reflect.ValueOf(vars)
	prefix, suffix := "${{", "}}"
	for i := 0; i < v.Elem().NumField(); i++ {
		currentValue := v.Elem().Field(i)
		if currentValue.Kind() != reflect.String {
			log.Println("skip variable because it is not string: ", currentValue.String())
			continue
		}
		currentStr := currentValue.Interface().(string)
		if strings.HasPrefix(currentStr, prefix) && strings.HasSuffix(currentStr, suffix) {
			buffer, _ := strings.CutPrefix(currentStr, prefix)
			buffer, _ = strings.CutSuffix(buffer, suffix)
			if envVar := os.Getenv(buffer); envVar != "" {
				v.Elem().Field(i).Set(reflect.ValueOf(envVar))
			} else {
				log.Fatalf("Unable to get variable %s from enviroment", buffer)
			}
		}
	}
}

func (v *Variables) IsValidForTelegram() bool {
	return v.TelegramBotToken != "" && v.TelegramChatID != ""
}

func (v *Variables) IsValidForSMTP() bool {
	return v.SmtpHost != "" && len(v.SmtpTo) != 0
}
