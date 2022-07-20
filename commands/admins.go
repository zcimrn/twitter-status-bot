package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zcimrn/twitter-status-bot/data"
	"github.com/zcimrn/twitter-status-bot/telegram"
	"github.com/zcimrn/twitter-status-bot/tools"
)

func getAdmins(chatId, messageId int) {
	admins := Data.GetAdmins()
	if len(admins) == 0 {
		telegram.SendMessage(chatId, "Админов пока нет", messageId)
		return
	}
	for i := 0; i < len(admins); i += 10 {
		text := "Админы:"
		for j := 0; j < 10 && i+j < len(admins); j++ {
			text += fmt.Sprintf("\n%d \\- `%d` \\- %s", i+j+1, admins[i+j].Id, tools.Escape(admins[i+j].Desc))
		}
		telegram.SendMessage(chatId, text, messageId)
	}
}

func addAdmin(chatId, messageId int, args []string) {
	if len(args) < 1 {
		telegram.SendMessage(chatId, "Не указан `id`", messageId)
		return
	}
	adminId, err := strconv.Atoi(args[0])
	if err != nil {
		telegram.SendMessage(chatId, "`id` должен быть числом", messageId)
		return
	}
	if len(args) < 2 {
		telegram.SendMessage(chatId, "Не указано описание", messageId)
		return
	}
	adminDesc := strings.Join(args[1:], " ")
	Data.AddAdmin(&data.Admin{adminId, adminDesc})
	telegram.SendMessage(chatId, fmt.Sprintf("Добавлен админ:\n`%d` %s", adminId, tools.Escape(adminDesc)), messageId)
}

func deleteAdmin(chatId, messageId int, args []string) {
	if len(args) < 1 {
		telegram.SendMessage(chatId, "Не указан `id`", messageId)
		return
	}
	adminId, err := strconv.Atoi(args[0])
	if err != nil {
		telegram.SendMessage(chatId, "`id` должен быть числом", messageId)
		return
	}
	Data.DeleteAdmin(adminId)
	telegram.SendMessage(chatId, fmt.Sprintf("Удалён админ c `id` `%d`", adminId), messageId)
}
