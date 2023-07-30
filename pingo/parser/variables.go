package parser

import (
    "strings"
    "reflect"
    "os"
    "log"
)

type Variables struct {
	TelegramBotToken string `yaml:"telegram_bot_token,omitempty"`
	TelegramChatID   string `yaml:"telegram_chat_id,omitempty"`
}

func (vars *Variables) PostParse() {
    // Check if there are any ${{...}} enviroment variables
    v := reflect.ValueOf(vars)
    prefix, suffix := "${{", "}}"
    for i:=0; i < v.Elem().NumField(); i++ {
        currentValue := v.Elem().Field(i).Interface().(string)
        if strings.HasPrefix(currentValue, prefix) && strings.HasSuffix(currentValue, suffix) {
            buffer, _ := strings.CutPrefix(currentValue, prefix)
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
