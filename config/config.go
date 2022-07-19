package config

import (
  "encoding/json"
  "fmt"
  "io"
  "log"
  "os"
  "sort"
  "sync"
)

type Admin struct {
  Id int `json:"id"`
  Desc string `json:"desc"`
}

type Chat struct {
  Id int `json:"id"`
  Desc string `json:"desc"`
}

type Config struct {
  FileName string `json:"file_name"`
  TelegramToken string `json:"telegram_token"`
  TwitterToken string `json:"twitter_token"`
  Admins []Admin `json:"admins"`
  Chats []Chat `json:"chats"`
  mutex sync.RWMutex `json:"-"`
}

func (config *Config) Init(fileName string) error {
  config.mutex.Lock()
  config.TelegramToken = ""
  config.TwitterToken = ""
  config.Admins = nil
  config.Chats = nil
  config.FileName = fileName
  err := config.load()
  if err != nil {
    config.mutex.Unlock()
    return err
  }
  if config.FileName == "" {
    config.mutex.Unlock()
    return fmt.Errorf("file_name is empty")
  }
  if config.TelegramToken == "" {
    config.mutex.Unlock()
    return fmt.Errorf("telegram_token is empty")
  }
  if config.TwitterToken == "" {
    config.mutex.Unlock()
    return fmt.Errorf("twitter_token is empty")
  }
  if len(config.Admins) == 0 {
    config.mutex.Unlock()
    return fmt.Errorf("not enough admins")
  }
  for i := 0; i < len(config.Admins); i++ {
    if config.Admins[i].Id == 0 {
      config.mutex.Unlock()
      return fmt.Errorf("admin %d has empty id", i + 1)
    }
    for j := i + 1; j < len(config.Admins); j++ {
      if config.Admins[j].Id == config.Admins[i].Id {
        config.mutex.Unlock()
        return fmt.Errorf("admins %d and %d had the same ids", i + 1, j + 1)
      }
    }
  }
  sort.Slice(config.Admins, func (i, j int) bool {
    return config.Admins[i].Id < config.Admins[j].Id
  })
  for i := 0; i < len(config.Chats); i++ {
    if config.Chats[i].Id == 0 {
      config.mutex.Unlock()
      return fmt.Errorf("chat %d has empty id", i + 1)
    }
    for j := i + 1; j < len(config.Chats); j++ {
      if config.Chats[j].Id == config.Chats[i].Id {
        config.mutex.Unlock()
        return fmt.Errorf("chats %d and %d had the same ids", i + 1, j + 1)
      }
    }
  }
  sort.Slice(config.Chats, func (i, j int) bool {
    return config.Chats[i].Id < config.Chats[j].Id
  })
  config.mutex.Unlock()
  return nil
}

func (config *Config) load() error {
  log.Printf("loading config from '%s'", config.FileName)
  file, err := os.Open(config.FileName)
  if err != nil {
    return err
  }
  bytes, err := io.ReadAll(file)
  file.Close()
  if err != nil {
    return err
  }
  err = json.Unmarshal(bytes, config)
  if err != nil {
    return err
  }
  log.Printf("config loaded from '%s'", config.FileName)
  return nil
}

func (config *Config) save() error {
  log.Printf("saving config to '%s'", config.FileName)
  bytes, err := json.MarshalIndent(config, "", "  ")
  if err != nil {
    return err
  }
  file, err := os.Create(config.FileName)
  if err != nil {
    return err
  }
  _, err = file.Write(bytes)
  if err != nil {
    return err
  }
  log.Printf("config saved to '%s'", config.FileName)
  return nil
}

func (config *Config) GetTelegramToken() string {
  config.mutex.RLock()
  telegramToken := config.TelegramToken
  config.mutex.RUnlock()
  return telegramToken
}

func (config *Config) SetTelegramToken(telegramToken string) {
  config.mutex.Lock()
  config.TelegramToken = telegramToken
  config.save()
  config.mutex.Unlock()
}

func (config *Config) GetTwitterToken() string {
  config.mutex.RLock()
  twitterToken := config.TwitterToken
  config.mutex.RUnlock()
  return twitterToken
}

func (config *Config) SetTwitterToken(twitterToken string) {
  config.mutex.Lock()
  config.TwitterToken = twitterToken
  config.save()
  config.mutex.Unlock()
}

func (config *Config) GetAdmins() []Admin {
  config.mutex.RLock()
  admins := append([]Admin{}, config.Admins...)
  config.mutex.RUnlock()
  return admins
}

func (config *Config) HasAdmin(id int) bool {
  config.mutex.RLock()
  i := sort.Search(len(config.Admins), func(i int) bool {
    return config.Admins[i].Id >= id
  })
  found := i < len(config.Admins) && config.Admins[i].Id == id
  config.mutex.RUnlock()
  return found
}

func (config *Config) AddAdmin(admin *Admin) bool {
  config.mutex.Lock()
  i := sort.Search(len(config.Admins), func (i int) bool {
    return config.Admins[i].Id >= admin.Id
  })
  if i < len(config.Admins) && config.Admins[i].Id == admin.Id {
    config.Admins[i] = *admin
    config.save()
    config.mutex.Unlock()
    return true
  }
  config.Admins = append(config.Admins[:i], append([]Admin{*admin}, config.Admins[i:]...)...)
  config.save()
  config.mutex.Unlock()
  return false
}

func (config *Config) DeleteAdmin(id int) bool {
  config.mutex.Lock()
  i := sort.Search(len(config.Admins), func (i int) bool {
    return config.Admins[i].Id >= id
  })
  if i == len(config.Admins) || config.Admins[i].Id != id {
    config.mutex.Unlock()
    return false
  }
  config.Admins = append(config.Admins[:i], config.Admins[i + 1:]...)
  config.save()
  config.mutex.Unlock()
  return true
}

func (config *Config) GetChats() []Chat {
  config.mutex.RLock()
  chats := append([]Chat{}, config.Chats...)
  config.mutex.RUnlock()
  return chats
}

func (config *Config) HasChat(id int) bool {
  config.mutex.RLock()
  i := sort.Search(len(config.Chats), func(i int) bool {
    return config.Chats[i].Id >= id
  })
  found := i < len(config.Chats) && config.Chats[i].Id == id
  config.mutex.RUnlock()
  return found
}

func (config *Config) AddChat(chat *Chat) bool {
  config.mutex.Lock()
  i := sort.Search(len(config.Chats), func(i int) bool {
    return config.Chats[i].Id >= chat.Id
  })
  if i < len(config.Chats) && config.Chats[i].Id == chat.Id {
    config.Chats[i] = *chat
    config.save()
    config.mutex.Unlock()
    return true
  }
  config.Chats = append(config.Chats[:i], append([]Chat{*chat}, config.Chats[i:]...)...)
  config.save()
  config.mutex.Unlock()
  return false
}

func (config *Config) DeleteChat(id int) bool {
  config.mutex.Lock()
  i := sort.Search(len(config.Chats), func(i int) bool {
    return config.Chats[i].Id >= id
  })
  if i == len(config.Chats) || config.Chats[i].Id != id {
    config.mutex.Unlock()
    return false
  }
  config.Chats = append(config.Chats[:i], config.Chats[i + 1:]...)
  config.save()
  config.mutex.Unlock()
  return true
}
