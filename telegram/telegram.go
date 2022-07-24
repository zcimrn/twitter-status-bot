package telegram

import (
	"encoding/json"
	"log"

	"github.com/zcimrn/twitter-status-bot/twitter"
)

type User struct {
	Id                      int    `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Username                string `json:"username"`
	LanguageCode            string `json:"language_code"`
	IsPremium               bool   `json:"is_premium"`
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
}

type Chat struct {
	Id int `json:"id"`
}

type Message struct {
	Id   int    `json:"message_id"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Update struct {
	Id      int     `json:"update_id"`
	Message Message `json:"message"`
}

func SendMessage(chatId int, text string, args ...int) error {
	replyToMessageId := 0
	if len(args) > 0 {
		replyToMessageId = args[0]
	}
	jsonReq := struct {
		ChatId                int    `json:"chat_id"`
		Text                  string `json:"text"`
		ParseMode             string `json:"parse_mode"`
		DisableWebPagePreview bool   `json:"disable_web_page_preview"`
		ReplyToMessageId      int    `json:"reply_to_message_id"`
	}{
		chatId,
		text,
		"MarkdownV2",
		true,
		replyToMessageId,
	}
	reqBody, err := json.Marshal(&jsonReq)
	if err != nil {
		return err
	}
	_, err = api("sendMessage", reqBody)
	if err != nil {
		return err
	}
	return nil
}

func GetUpdates(offset int) ([]Update, error) {
	jsonReq := struct {
		Offset         int      `json:"offset"`
		Timeout        int      `json:"timeout"`
		AllowedUpdates []string `json:"allowed_updates"`
	}{
		offset,
		1<<31 - 1,
		[]string{"message"},
	}
	reqBody, err := json.Marshal(&jsonReq)
	if err != nil {
		return nil, err
	}
	respBody, err := api("getUpdates", reqBody)
	if err != nil {
		return nil, err
	}
	var jsonResp struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}
	err = json.Unmarshal(respBody, &jsonResp)
	if err != nil {
		return nil, err
	}
	return jsonResp.Result, err
}

func SendUpdates(user *twitter.User, updates []twitter.User) {
	log.Printf("[%s] sending updates...", user.Username)
	for i := 0; i < len(updates); i += 10 {
		text := user.Markdown() + " подписался на:"
		for j := 0; j < 10 && i+j < len(updates); j++ {
			text += "\n" + updates[i+j].Markdown()
		}
		for _, id := range user.GetChatIds() {
			err := SendMessage(id, text)
			if err != nil {
				log.Printf("[%s] chat %d send message error: '%s'", user.Username, id, err)
			}
		}
	}
	log.Printf("[%s] updates sent", user.Username)
}
