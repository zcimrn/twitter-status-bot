package data

import (
	"github.com/zcimrn/twitter-status-bot/telegram"
)

func (data *Data) GetTelegramToken() string {
	data.mutex.RLock()
	telegramToken := data.TelegramToken
	data.mutex.RUnlock()
	return telegramToken
}

func (data *Data) SetTelegramToken(telegramToken string) {
	data.mutex.Lock()
	data.TelegramToken = telegramToken
	telegram.SetToken(telegramToken)
	data.save()
	data.mutex.Unlock()
}
