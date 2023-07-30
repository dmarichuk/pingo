package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var TgClient *http.Client

func init() {
	TgClient = &http.Client{}
}

type TelegramTask struct {
	BotToken     string
	ChatID       string
	AlertMessage string
}

func NewTelegramTask(token, chatID, alertMessage string) TelegramTask {
	return TelegramTask{
		BotToken:     token,
		ChatID:       chatID,
		AlertMessage: alertMessage,
	}
}

func (task TelegramTask) Launch() {
	data := map[string]interface{}{
		"chat_id": task.ChatID,
		"text":    task.AlertMessage,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Print("Error marshaling JSON:", err)
		return
	}

	telegramEndpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", task.BotToken)
	req, err := http.NewRequest("POST", telegramEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Print("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Print("Error sending telegram message:", resp.StatusCode)
	}

}
