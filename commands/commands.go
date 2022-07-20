package commands

import (
  "log"
  "strings"

  "github.com/zcimrn/twitter-status-bot/data"
  "github.com/zcimrn/twitter-status-bot/telegram"
)

var Data *data.Data

func Exec(message *telegram.Message) {
  log.Printf("message: '%+v'", message)
  if len(message.Text) == 0 {
    log.Printf("error: 'message without text'")
    return
  }
  if !Data.HasAdmin(message.From.Id) && !Data.HasAdmin(message.Chat.Id) {
    log.Printf("error: 'permissions error'")
    return
  }
  chatId := message.Chat.Id
  messageId := message.Id
  args := strings.Fields(message.Text)
  command, _, _ := strings.Cut(strings.ToLower(args[0]), "@")
  args = args[1:]
  switch command {
  case "/help":
    help(chatId, messageId)
  case "/get_telegram_token":
    getTelegramToken(chatId, messageId)
  case "/set_telegram_token":
    setTelegramToken(chatId, messageId, args)
  case "/get_twitter_token":
    getTwitterToken(chatId, messageId)
  case "/set_twitter_token":
    setTwitterToken(chatId, messageId, args)
  case "/get_admins":
    getAdmins(chatId, messageId)
  case "/add_admin":
    addAdmin(chatId, messageId, args)
  case "/delete_admin":
    deleteAdmin(chatId, messageId, args)
  case "/get_chats":
    getChats(chatId, messageId)
  case "/add_chat":
    addChat(chatId, messageId, args)
  case "/delete_chat":
    deleteChat(chatId, messageId, args)
  case "/get_users":
    getUsers(chatId, messageId, args)
  case "/add_users":
    addUsers(chatId, messageId, args)
  case "/delete_users":
    deleteUsers(chatId, messageId, args)
  default:
    unknown(chatId, messageId, command)
  }
}
