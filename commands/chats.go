package commands

import (
  "fmt"
  "strconv"
  "strings"

  "github.com/zcimrn/twitter-status-bot/config"
  "github.com/zcimrn/twitter-status-bot/telegram"
  "github.com/zcimrn/twitter-status-bot/tools"
)

func getChats(chatId, messageId int) {
  chats := Config.GetChats()
  if len(chats) == 0 {
    telegram.SendMessage(chatId, "Чатов пока нет", messageId)
    return
  }
  for i := 0; i < len(chats); i += 10 {
    text := "Чаты:"
    for j := 0; j < 10 && i + j < len(chats); j++ {
      text += fmt.Sprintf("\n%d \\- `%d` \\- %s", i + j + 1, chats[i + j].Id, tools.Escape(chats[i + j].Desc))
    }
    telegram.SendMessage(chatId, text, messageId)
  }
}

func addChat(chatId, messageId int, args []string) {
  if len(args) < 1 {
    telegram.SendMessage(chatId, "Не указан `id`", messageId)
    return
  }
  commandChatId, err := strconv.Atoi(args[0])
  if err != nil {
    telegram.SendMessage(chatId, "`id` должен быть числом", messageId)
    return
  }
  if len(args) < 2 {
    telegram.SendMessage(chatId, "Не указано описание", messageId)
    return
  }
  commandChatDesc := strings.Join(args[1:], " ")
  Config.AddChat(&config.Chat{ commandChatId, commandChatDesc })
  telegram.SendMessage(chatId, fmt.Sprintf("Добавлен чат\n`%d` %s", commandChatId, tools.Escape(commandChatDesc)), messageId)
}

func deleteChat(chatId, messageId int, args []string) {
  if len(args) < 1 {
    telegram.SendMessage(chatId, "Не указан `id`", messageId)
    return
  }
  commandChatId, err := strconv.Atoi(args[0])
  if err != nil {
    telegram.SendMessage(chatId, "`id` должен быть числом", messageId)
    return
  }
  if Config.DeleteChat(commandChatId) {
    telegram.SendMessage(chatId, fmt.Sprintf("Удалён чат c `id` `%d`", commandChatId), messageId)
  } else {
    telegram.SendMessage(chatId, fmt.Sprintf("Нет чата с `id` `%d`", commandChatId), messageId)
  }
}


