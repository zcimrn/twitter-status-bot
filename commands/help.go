package commands

import (
  "fmt"

  "github.com/zcimrn/twitter-status-bot/telegram"
  "github.com/zcimrn/twitter-status-bot/tools"
)

func help(chatId, messageId int) {
  telegram.SendMessage(
    chatId,
    "Справка:\n`/help`\n\n" +
    "Получить Telegram token:\n`/get_telegram_token`\n\n" +
    "Установить Telegram token:\n`/set_telegram_token <telegram_token>`\n\n" +
    "Получить Twitter token:\n`/get_twitter_token`\n\n" +
    "Установить Twitter token:\n`/set_twitter_token <twitter_token>`\n\n" +
    "Получить список админов:\n`/get_admins`\n\n" +
    "Добавить админа:\n`/add_admin <admin_id> <description>`\n\n" +
    "Удалить админа:\n`/delete_admin <admin_id>`\n\n" +
    "Получить список чатов:\n`/get_chats`\n\n" +
    "Добавить чат:\n`/add_chat <chat_id> <description>`\n\n" +
    "Удалить чат:\n`/detele_chat <chat_id>`\n\n" +
    "Получить список аккаунтов:\n`/get_users`\n\n" +
    "Добавить аккаунты:\n`/add_users <chat_id> <usernames...>`\n\n" +
    "Удалить аккаунты:\n`/delete_users <chat_id> <usernames...>`",
    messageId,
  )
}

func unknown(chatId, messageId int, command string) {
  telegram.SendMessage(chatId, fmt.Sprintf("Команды `%s` нет\nПопробуйте /help", tools.EscapeCode(command)), messageId)
}
