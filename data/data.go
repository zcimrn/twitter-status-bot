package data

import (
  "encoding/json"
  "fmt"
  "io"
  "log"
  "os"
  "sort"
  "sync"

  "github.com/zcimrn/twitter-status-bot/twitter"
)

type Data struct {
  FileName string `json:"file_name"`
  LastUserIndex int `json:"last_user_index"`
  Users []twitter.User `json:"users"`
  mutex sync.RWMutex `json:"-"`
  emptyMutex sync.Mutex `json:"-"`
}

func (data *Data) Init(fileName string) error {
  data.mutex.Lock()
  data.emptyMutex.Lock()
  data.LastUserIndex = -1
  data.Users = nil
  data.FileName = fileName
  err := data.load()
  if err != nil {
    log.Printf("error: '%s'", err)
  }
  if data.FileName == "" {
    data.emptyMutex.Unlock()
    data.mutex.Unlock()
    return fmt.Errorf("file_name is empty")
  }
  if data.LastUserIndex < -1 {
    data.emptyMutex.Unlock()
    data.mutex.Unlock()
    return fmt.Errorf("last_user_index < -1")
  }
  if data.LastUserIndex >= len(data.Users) {
    data.emptyMutex.Unlock()
    data.mutex.Unlock()
    return fmt.Errorf("last_user_index > %d", len(data.Users) - 1)
  }
  for i := 0; i < len(data.Users); i++ {
    err = data.Users[i].Validate()
    if err != nil {
      data.emptyMutex.Unlock()
      data.mutex.Unlock()
      return fmt.Errorf("user %d error: '%s'", i + 1, err)
    }
    for j := i + 1; j < len(data.Users); j++ {
      if data.Users[j].Id == data.Users[i].Id {
        data.emptyMutex.Unlock()
        data.mutex.Unlock()
        return fmt.Errorf("users %d and %d has the same ids", i + 1, j + 1)
      }
    }
  }
  sort.Slice(data.Users, func (i, j int) bool {
    return data.Users[i].Id < data.Users[j].Id
  })
  if len(data.Users) > 0 {
    data.emptyMutex.Unlock()
  }
  data.mutex.Unlock()
  return nil
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
  bytes, err := json.MarshalIndent(data, "", "  ")
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

func (data *Data) GetUsers() []twitter.User {
  data.mutex.RLock()
  var users []twitter.User
  for i := 0; i < len(data.Users); i++ {
    users = append(users, *data.Users[i].Copy())
  }
  data.mutex.RUnlock()
  return users
}

func (data *Data) HasUser(id string) bool {
  data.mutex.RLock()
  i := sort.Search(len(data.Users), func(i int) bool {
    return data.Users[i].Id >= id
  })
  found := i < len(data.Users) && data.Users[i].Id == id
  data.mutex.RUnlock()
  return found
}

func (data *Data) GetUsersByChatId(id int) []twitter.User {
  var users []twitter.User
  data.mutex.RLock()
  for i := 0; i < len(data.Users); i++ {
    if data.Users[i].HasChatId(id) {
      users = append(users, *data.Users[i].Copy())
    }
  }
  data.mutex.RUnlock()
  return users
}

func (data *Data) UpdateUser(user *twitter.User) bool {
  data.mutex.Lock()
  i := sort.Search(len(data.Users), func (i int) bool {
    return data.Users[i].Id >= user.Id
  })
  if i == len(data.Users) || data.Users[i].Id != user.Id {
    data.mutex.Unlock()
    return false
  }
  data.Users[i] = *user
  data.save()
  data.mutex.Unlock()
  return true
}

func (data *Data) AddUser(user *twitter.User) bool {
  data.mutex.Lock()
  i := sort.Search(len(data.Users), func (i int) bool {
    return data.Users[i].Id >= user.Id
  })
  if i < len(data.Users) && data.Users[i].Id == user.Id {
    data.emptyMutex.Lock()
    for _, chatId := range user.GetChatIds() {
      data.Users[i].AddChatId(chatId)
    }
    data.save()
    data.emptyMutex.Unlock()
    data.mutex.Unlock()
    return true
  }
  if len(data.Users) > 0 {
    data.emptyMutex.Lock()
  }
  data.Users = append(data.Users[:i], append([]twitter.User{*user}, data.Users[i:]...)...)
  if data.LastUserIndex >= i {
    data.LastUserIndex++
  }
  data.save()
  data.emptyMutex.Unlock()
  data.mutex.Unlock()
  return false
}

func (data *Data) DeleteUser(user *twitter.User) bool {
  data.mutex.Lock()
  i := sort.Search(len(data.Users), func (i int) bool {
    return data.Users[i].Id >= user.Id
  })
  if i == len(data.Users) || data.Users[i].Id != user.Id {
    data.mutex.Unlock()
    return false
  }
  data.emptyMutex.Lock()
  for _, id := range user.GetChatIds() {
    data.Users[i].DeleteChatId(id)
  }
  if len(data.Users[i].GetChatIds()) == 0 {
    data.Users = append(data.Users[:i], data.Users[i + 1:]...)
    if data.LastUserIndex >= i {
      data.LastUserIndex--
    }
  }
  data.save()
  if len(data.Users) > 0 {
    data.emptyMutex.Unlock()
  }
  data.mutex.Unlock()
  return true
}

func (data *Data) GetNextUser() *twitter.User {
  data.emptyMutex.Lock()
  data.LastUserIndex = (data.LastUserIndex + 1) % len(data.Users)
  user := data.Users[data.LastUserIndex].Copy()
  data.save()
  data.emptyMutex.Unlock()
  return user
}
