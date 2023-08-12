package parser

import (
	"log"
	"os"
	"reflect"
	"strings"
)

type Settings struct {
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

func (s *Settings) PostParse() {
	// Check if there are any ${{...}} enviroment variables
	v := reflect.ValueOf(s)
	prefix, suffix := "${{", "}}"
	for i := 0; i < v.Elem().NumField(); i++ {
		currentValue := v.Elem().Field(i)
		if currentValue.Kind() != reflect.String {
			continue
		}
		currentStr := currentValue.Interface().(string)
		if strings.HasPrefix(currentStr, prefix) && strings.HasSuffix(currentStr, suffix) {
			buffer, _ := strings.CutPrefix(currentStr, prefix)
			buffer, _ = strings.CutSuffix(buffer, suffix)
			buffer = strings.Trim(buffer, " ")
			if envVar := os.Getenv(buffer); envVar != "" {
				v.Elem().Field(i).Set(reflect.ValueOf(envVar))
			} else {
				log.Fatalf("Unable to get variable %s from enviroment", buffer)
			}
		}
	}
}

func (s *Settings) IsValidForTelegram() bool {
	return s.TelegramBotToken != "" && s.TelegramChatID != ""
}

func (s *Settings) IsValidForSMTP() bool {
	return s.SmtpHost != "" && len(s.SmtpTo) != 0
}
