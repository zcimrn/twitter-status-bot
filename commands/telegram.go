package commands

import (
  "fmt"

  "github.com/zcimrn/twitter-status-bot/telegram"
  "github.com/zcimrn/twitter-status-bot/tools"
)

func getTelegramToken(chatId, messageId int) {
  telegram.SendMessage(chatId, fmt.Sprintf("Telegram token:\n`%s`", tools.EscapeCode(Config.GetTelegramToken())), messageId)
}

func setTelegramToken(chatId, messageId int, args []string) {
  if len(args) == 0 {
    telegram.SendMessage(chatId, "Не указан `telegram_token`", messageId)
    return
  }
  token := args[0]
  if !telegram.TestToken(token) {
    telegram.SendMessage(chatId, fmt.Sprintf("Не удалось установить Telegram token:\n`%s`", tools.EscapeCode(token)), messageId)
    return
  }
  Config.SetTelegramToken(token)
  telegram.SendMessage(chatId, fmt.Sprintf("Установлен Telegram token:\n`%s`", tools.EscapeCode(token)), messageId)
}
