package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pingo/job"
	"time"
)

var TgClient *http.Client

func init() {
	TgClient = &http.Client{}
}

type TelegramAlert struct {
	BotToken string
	ChatID   string
}

func NewTelegramAlert(token, chatID string) *TelegramAlert {
	return &TelegramAlert{
		BotToken: token,
		ChatID:   chatID,
	}
}

func (a *TelegramAlert) Send(msg []byte) {
	data := struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: a.ChatID,
		Text:   string(msg),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Print("[ERR] Error marshaling JSON:", err)
		return
	}

	telegramEndpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", a.BotToken)
	req, err := http.NewRequest("POST", telegramEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Print("[ERR] Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("[ERR] Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Print("[ERR] Error sending telegram message:", resp.StatusCode)
	}

}

func (a *TelegramAlert) GenerateMessage(j *job.Job) []byte {
	return []byte(fmt.Sprintf("Job %s; Status %s; Time %s", j.Name, j.Status, j.TS.Format(time.RFC3339)))
}
