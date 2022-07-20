package commands

import (
  "fmt"
  "log"
  "sort"
  "strconv"
  "strings"

  "github.com/zcimrn/twitter-status-bot/telegram"
  "github.com/zcimrn/twitter-status-bot/tools"
  "github.com/zcimrn/twitter-status-bot/twitter"
)

func getUsers(chatId, messageId int, args []string) {
  if len(args) == 0 {
    telegram.SendMessage(chatId, "Не указан `chat_id`", messageId)
    return
  }
  commandChatId, err := strconv.Atoi(args[0])
  if err != nil {
    telegram.SendMessage(chatId, "`chat_id` должен быть числом", messageId)
    return
  }
  if !Data.HasChat(commandChatId) {
    telegram.SendMessage(chatId, fmt.Sprintf("Чата `%d` нет в списке", commandChatId), messageId)
    return
  }
  users := Data.GetUsersByChatId(commandChatId)
  if len(users) == 0 {
    telegram.SendMessage(chatId, fmt.Sprintf("Для чата `%d` пока нет аккаунтов", commandChatId), messageId)
    return
  }
  sort.Slice(users, func (i, j int) bool {
    return users[i].Username < users[j].Username
  })
  for i := 0; i < len(users); i += 50 {
    text := fmt.Sprintf("Аккаунты для чата `%d`:", commandChatId)
    for j := 0; j < 50 && i + j < len(users); j++ {
      text += fmt.Sprintf("\n%d \\- `%s`", i + j + 1, tools.EscapeCode(users[i + j].Username))
    }
    telegram.SendMessage(chatId, text, messageId)
  }
}

func addUsers(chatId, messageId int, args []string) {
  if len(args) == 0 {
    telegram.SendMessage(chatId, "Не указан `chat_id`", messageId)
    return
  }
  commandChatId, err := strconv.Atoi(args[0])
  if err != nil {
    telegram.SendMessage(chatId, "`chat_id` должен быть числом", messageId)
    return
  }
  if !Data.HasChat(commandChatId) {
    telegram.SendMessage(chatId, fmt.Sprintf("Чата `%d` нет в списке", commandChatId), messageId)
    return
  }
  usernames := args[1:]
  if len(usernames) == 0 {
    telegram.SendMessage(chatId, "Нужно указать хотя бы один аккаунт", messageId)
  }
  if len(usernames) > 250 {
    telegram.SendMessage(chatId, "Не получится добавить более 250 аккаунтов за раз", messageId)
    return
  }
  for i := 0; i < len(usernames); i++ {
    usernames[i] = strings.ToLower(usernames[i])
  }
  sort.Strings(usernames)
  var okIndexes, errorIndexes []int
  for i := 0; i < len(usernames); i++ {
    user, err := twitter.GetUserByUsername(usernames[i])
    if err != nil {
      log.Printf("error: '%s'", err)
      errorIndexes = append(errorIndexes, i)
      continue
    }
    user.AddChatId(commandChatId)
    Data.AddUser(user)
    okIndexes = append(okIndexes, i)
  }
  if len(okIndexes) > 0 {
    text := fmt.Sprintf("Для чата `%d` добавлены аккаунты:```\n", commandChatId)
    for _, i := range okIndexes {
      text += tools.EscapeCode(usernames[i]) + "\n"
    }
    text += "```"
    telegram.SendMessage(chatId, text, messageId)
  }
  if len(errorIndexes) > 0 {
    text := fmt.Sprintf("Для чата `%d` не получилось добавить аккаунты:```\n", commandChatId)
    for _, i := range errorIndexes {
      text += tools.EscapeCode(usernames[i]) + "\n"
    }
    text += "```"
    telegram.SendMessage(chatId, text, messageId)
  }
}

func deleteUsers(chatId, messageId int, args []string) {
  if len(args) == 0 {
    telegram.SendMessage(chatId, "Не указан `chat_id`", messageId)
    return
  }
  commandChatId, err := strconv.Atoi(args[0])
  if err != nil {
    telegram.SendMessage(chatId, "`chat_id` должен быть числом", messageId)
    return
  }
  if !Data.HasChat(commandChatId) {
    telegram.SendMessage(chatId, fmt.Sprintf("Чата `%d` нет в списке", commandChatId), messageId)
    return
  }
  usernames := args[1:]
  if len(usernames) == 0 {
    telegram.SendMessage(chatId, "Нужно указать хотя бы один аккаунт", messageId)
  }
  if len(usernames) > 250 {
    telegram.SendMessage(chatId, "Не получится удалить более 250 аккаунтов за раз", messageId)
    return
  }
  for i := 0; i < len(usernames); i++ {
    usernames[i] = strings.ToLower(usernames[i])
  }
  sort.Strings(usernames)
  var okIndexes, errorIndexes []int
  for i := 0; i < len(usernames); i++ {
    user, err := twitter.GetUserByUsername(usernames[i])
    if err != nil {
      log.Printf("error: '%s'", err)
      errorIndexes = append(errorIndexes, i)
      continue
    }
    user.AddChatId(commandChatId)
    Data.DeleteUser(user)
    okIndexes = append(okIndexes, i)
  }
  if len(okIndexes) > 0 {
    text := fmt.Sprintf("Для чата `%d` удалены аккаунты:```\n", commandChatId)
    for _, i := range okIndexes {
      text += tools.EscapeCode(usernames[i]) + "\n"
    }
    text += "```"
    telegram.SendMessage(chatId, text, messageId)
  }
  if len(errorIndexes) > 0 {
    text := fmt.Sprintf("Для чата `%d` не получилось удалить аккаунты:```\n", commandChatId)
    for _, i := range errorIndexes {
      text += tools.EscapeCode(usernames[i]) + "\n"
    }
    text += "```"
    telegram.SendMessage(chatId, text, messageId)
  }
}


