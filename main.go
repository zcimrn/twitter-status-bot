package main

import (
  "log"
  "sync"
  "time"

  "github.com/zcimrn/twitter-status-bot/commands"
  "github.com/zcimrn/twitter-status-bot/config"
  "github.com/zcimrn/twitter-status-bot/data"
  "github.com/zcimrn/twitter-status-bot/telegram"
  "github.com/zcimrn/twitter-status-bot/twitter"
)

var (
  Config = &config.Config{}
  Data = &data.Data{}
)

func monitorTwitter(delay time.Duration) {
  for {
    user := Data.GetNextUser()
    updates := user.Update(delay)
    user.Pretty()
    Data.UpdateUser(user)
    if len(updates) > 0 {
      telegram.SendUpdates(user, updates)
    }
  }
}

func monitorTelegram() {
  lastUpdateId := 0
  for {
    updates, err := telegram.GetUpdates(lastUpdateId + 1)
    if err != nil {
      log.Printf("error: '%s'", err)
    }
    for i := 0; i < len(updates); i++ {
      go commands.Exec(&updates[i].Message)
      lastUpdateId = updates[i].Id
    }
  }
}

func main() {
  log.Printf("initializing...")
  err := Config.Init("config.json")
  if err != nil {
    log.Printf("error: '%s'", err)
    return
  }
  log.Printf("config: '%+v'", Config)
  telegram.Config = Config
  twitter.Config = Config
  commands.Config = Config
  err = Data.Init("data.json")
  if err != nil {
    log.Printf("error: '%s'", err)
    return
  }
  commands.Data = Data
  log.Printf("initialized")
  var waitGroup sync.WaitGroup
  waitGroup.Add(2)
  go monitorTwitter(60 * time.Second)
  go monitorTelegram()
  waitGroup.Wait()
}
