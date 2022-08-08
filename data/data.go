package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/zcimrn/twitter-status-bot/telegram"
	"github.com/zcimrn/twitter-status-bot/twitter"
)

type Admin struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
}

type Chat struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
}

type Data struct {
	FileName      string         `json:"file_name"`
	TelegramToken string         `json:"telegram_token"`
	TwitterToken  string         `json:"twitter_token"`
	Admins        []Admin        `json:"admins"`
	Chats         []Chat         `json:"chats"`
	LastIndex     int            `json:"last_user_index"`
	Users         []twitter.User `json:"users"`
	mutex         sync.RWMutex   `json:"-"`
}

func (data *Data) load() error {
	log.Printf("loading data from '%s'", data.FileName)
	file, err := os.Open(data.FileName)
	if err != nil {
		return err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, data)
	if err != nil {
		return err
	}
	log.Printf("data loaded from '%s'", data.FileName)
	return nil
}

func (data *Data) save() error {
	log.Printf("saving data to '%s'", data.FileName)
	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	file, err := os.Create(data.FileName)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	file.Close()
	if err != nil {
		return err
	}
	log.Printf("data saved to '%s'", data.FileName)
	return nil
}

func (data *Data) Init(fileName string) error {
	data.mutex.Lock()
	data.FileName = fileName
	data.TelegramToken = ""
	data.TwitterToken = ""
	data.Admins = nil
	data.Chats = nil
	data.LastIndex = -1
	data.Users = nil
	if err := data.load(); err != nil {
		return err
	}
	if data.FileName == "" {
		data.mutex.Unlock()
		return fmt.Errorf("file_name is empty")
	}
	if data.TelegramToken == "" {
		data.mutex.Unlock()
		return fmt.Errorf("telegram_token is empty")
	}
	if err := telegram.TestToken(data.TelegramToken); err != nil {
		data.mutex.Unlock()
		return err
	}
	telegram.SetToken(data.TelegramToken)
	if data.TwitterToken == "" {
		data.mutex.Unlock()
		return fmt.Errorf("twitter_token is empty")
	}
	if !twitter.TestToken(data.TwitterToken) {
		data.mutex.Unlock()
		return fmt.Errorf("twitter_token doesn't work")
	}
	twitter.SetToken(data.TwitterToken)
	if len(data.Admins) == 0 {
		data.mutex.Unlock()
		return fmt.Errorf("not enough admins")
	}
	for i := 0; i < len(data.Admins); i++ {
		for j := i + 1; j < len(data.Admins); j++ {
			if data.Admins[j].Id == data.Admins[i].Id {
				data.mutex.Unlock()
				return fmt.Errorf("admins %d and %d have the same ids", i+1, j+1)
			}
		}
	}
	sort.Slice(data.Admins, func(i, j int) bool {
		return data.Admins[i].Id < data.Admins[j].Id
	})
	for i := 0; i < len(data.Chats); i++ {
		for j := i + 1; j < len(data.Chats); j++ {
			if data.Chats[j].Id == data.Chats[i].Id {
				data.mutex.Unlock()
				return fmt.Errorf("chats %d and %d have the same ids", i+1, j+1)
			}
		}
	}
	sort.Slice(data.Chats, func(i, j int) bool {
		return data.Chats[i].Id < data.Chats[j].Id
	})
	if data.LastIndex < -1 {
		data.mutex.Unlock()
		return fmt.Errorf("last_index < -1")
	}
	if data.LastIndex >= len(data.Users) {
		data.mutex.Unlock()
		return fmt.Errorf("last_index > %d", len(data.Users)-1)
	}
	for i := 0; i < len(data.Users); i++ {
		if err := data.Users[i].Validate(); err != nil {
			data.mutex.Unlock()
			return fmt.Errorf("user %d error: '%s'", i+1, err)
		}
		for j := i + 1; j < len(data.Users); j++ {
			if data.Users[j].Id == data.Users[i].Id {
				data.mutex.Unlock()
				return fmt.Errorf("users %d and %d has the same ids", i+1, j+1)
			}
		}
	}
	sort.Slice(data.Users, func(i, j int) bool {
		return data.Users[i].Id < data.Users[j].Id
	})
	data.save()
	data.mutex.Unlock()
	return nil
}
