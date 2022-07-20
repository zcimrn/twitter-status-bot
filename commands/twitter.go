package commands

import (
	"fmt"

	"github.com/zcimrn/twitter-status-bot/telegram"
	"github.com/zcimrn/twitter-status-bot/tools"
	"github.com/zcimrn/twitter-status-bot/twitter"
)

func getTwitterToken(chatId, messageId int) {
	telegram.SendMessage(chatId, fmt.Sprintf("Twitter token:\n`%s`", tools.EscapeCode(Data.GetTwitterToken())), messageId)
}

func setTwitterToken(chatId, messageId int, args []string) {
	if len(args) == 0 {
		telegram.SendMessage(chatId, "Не указан `twitter_token`", messageId)
		return
	}
	token := args[0]
	if !twitter.TestToken(token) {
		telegram.SendMessage(chatId, fmt.Sprintf("Не удалось установить Twitter token:\n`%s`", tools.EscapeCode(token)), messageId)
		return
	}
	Data.SetTwitterToken(token)
	telegram.SendMessage(chatId, fmt.Sprintf("Установлен Twitter token:\n`%s`", tools.EscapeCode(token)), messageId)
}
