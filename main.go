package main

import (
	"log"
	"sync"
	"time"

	"github.com/zcimrn/twitter-status-bot/commands"
	"github.com/zcimrn/twitter-status-bot/data"
	"github.com/zcimrn/twitter-status-bot/telegram"
	"github.com/zcimrn/twitter-status-bot/twitter"
)

func monitorTwitter(data *data.Data, delay time.Duration) {
	for {
		var user *twitter.User
		for {
			user = data.GetNextUser()
			if user != nil {
				break
			}
			log.Printf("error: 'empty data'")
			log.Printf("waiting...")
			time.Sleep(delay)
		}
		updates := user.Update(delay)
		if data.UpdateUser(user) && len(updates) > 0 {
			telegram.SendUpdates(user, updates)
		}
	}
}

func monitorTelegram() {
	lastUpdateId := 0
	for {
		updates, err := telegram.GetUpdates(lastUpdateId + 1)
		if err != nil {
			log.Printf("telegram get updates error: '%s'", err)
		}
		for i := 0; i < len(updates); i++ {
			go commands.Exec(&updates[i].Message)
			lastUpdateId = updates[i].Id
		}
	}
}

func main() {
	log.Printf("initializing...")
	data := &data.Data{}
	if err := data.Init("data.json"); err != nil {
		log.Printf("error: '%s'", err)
		return
	}
	commands.Data = data
	log.Printf("initialized")
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go monitorTwitter(data, 60*time.Second)
	go monitorTelegram()
	waitGroup.Wait()
}
