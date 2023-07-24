package parser

// TODO post process variables after reading from file and get env variables for string with $ prefix
type Variables struct {
	TelegramBotToken string `yaml:"telegram_bot_token,omitempty"`
	TelegramChatID   string `yaml:"telegram_chat_id,omitempty"`
}

func (v *Variables) IsValidForTelegram() bool {
	return v.TelegramBotToken != "" && v.TelegramChatID != ""
}
