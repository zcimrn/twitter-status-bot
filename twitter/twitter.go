package twitter

import (
  "encoding/json"
  "fmt"
  "log"
  "sort"
  "strconv"
  "strings"
  "time"

  "github.com/zcimrn/twitter-status-bot/tools"
)

type User struct {
  Id string `json:"id"`
  Username string `json:"username"`
  Name string `json:"name"`
  Initialized bool `json:"initialized"`
  Followings []User `json:"followings"`
  ChatIds []int `json:"chat_ids"`
}

func getUserById(id string) (*User, error) {
  respBody, err := api("https://api.twitter.com/2/users/" + id)
  if err != nil {
    return nil, err
  }
  var jsonResp struct {
    Data User `json:"data"`
  }
  err = json.Unmarshal(respBody, &jsonResp)
  if err != nil {
    return nil, err
  }
  jsonResp.Data.Username = strings.ToLower(jsonResp.Data.Username)
  return &jsonResp.Data, nil
}

func GetUserByUsername(username string) (*User, error) {
  return getUserById("by/username/" + username)
}

func getFollowings(userId string, delay time.Duration) ([]User, error) {
  var users []User
  for nextToken := "";; {
    time.Sleep(delay)
    query := "https://api.twitter.com/2/users/" + userId + "/following?max_results=1000"
    if nextToken != "" {
      query += "&pagination_token=" + nextToken
    }
    respBody, err := api(query)
    if err != nil {
      return nil, err
    }
    var jsonResp struct {
      Data []User `json:"data"`
      Meta struct {
        ResultCount int `json:"result_count"`
        NextToken string `json:"next_token"`
      } `json:"meta"`
    }
    err = json.Unmarshal(respBody, &jsonResp)
    if err != nil {
      return nil, err
    }
    for i := 0; i < len(jsonResp.Data); i++ {
      jsonResp.Data[i].Username = strings.ToLower(jsonResp.Data[i].Username)
      users = append(users, jsonResp.Data[i])
    }
    nextToken = jsonResp.Meta.NextToken
    if nextToken == "" {
      break
    }
  }
  return users, nil
}

func (user *User) Update(delay time.Duration) []User {
  log.Printf("[%s] updating...", user.Username)
  newUser, err := getUserById(user.Id)
  if err != nil {
    log.Printf("[%s] error: '%s'", user.Username, err)
    return nil
  }
  user.Username = newUser.Username
  user.Name = newUser.Name
  log.Printf("[%s] getting followings...", user.Username)
  newUser.Followings, err = getFollowings(user.Id, delay)
  if err != nil {
    log.Printf("[%s] error: '%s'", user.Username, err)
    return nil
  }
  log.Printf("[%s] got %d followings", user.Username, len(newUser.Followings))
  var newFollowings []User
  if user.Initialized {
    log.Printf("[%s] checking difference...", user.Username)
    for i := 0; i < len(newUser.Followings); i++ {
      found := false
      for j := 0; j < len(user.Followings); j++ {
        if newUser.Followings[i].Id == user.Followings[j].Id {
          found = true
          break
        }
      }
      if !found {
        newFollowings = append(newFollowings, newUser.Followings[i])
      }
    }
    log.Printf("[%s] found %d new followings", user.Username, len(newFollowings))
  } else {
    log.Printf("[%s] not initialized", user.Username)
    user.Initialized = true
  }
  user.Followings = newUser.Followings
  log.Printf("[%s] updated", user.Username)
  return newFollowings
}

func (user *User) GetChatIds() []int {
  return append([]int{}, user.ChatIds...)
}

func (user *User) HasChatId(id int) bool {
  i := sort.Search(len(user.ChatIds), func(i int) bool {
    return user.ChatIds[i] >= id
  })
  return i < len(user.ChatIds) && user.ChatIds[i] == id
}

func (user *User) AddChatId(id int) bool {
  i := sort.Search(len(user.ChatIds), func(i int) bool {
    return user.ChatIds[i] >= id
  })
  if i < len(user.ChatIds) && user.ChatIds[i] == id {
    return true
  }
  user.ChatIds = append(user.ChatIds[:i], append([]int{id}, user.ChatIds[i:]...)...)
  return false
}

func (user *User) DeleteChatId(id int) bool {
  i := sort.Search(len(user.ChatIds), func (i int) bool {
    return user.ChatIds[i] >= id
  })
  if i == len(user.ChatIds) || user.ChatIds[i] != id {
    return false
  }
  user.ChatIds = append(user.ChatIds[:i], user.ChatIds[i + 1:]...)
  return true
}

func (user *User) Pretty() {
  user.Username = strings.ToLower(user.Username)
  sort.Slice(user.Followings, func (i, j int) bool {
    return user.Followings[i].Id < user.Followings[j].Id
  })
  sort.Ints(user.ChatIds)
}

func (user *User) Validate() error {
  if user.Id == "" {
    return fmt.Errorf("user has empty id")
  }
  _, err := strconv.Atoi(user.Id)
  if err != nil {
    return fmt.Errorf("user has non int id")
  }
  for i := 0; i < len(user.Followings); i++ {
    err = user.Followings[i].Validate()
    if err != nil {
      return fmt.Errorf("following %d error: '%s'", i + 1, err)
    }
    for j := i + 1; j < len(user.Followings); j++ {
      if user.Followings[j].Id == user.Followings[i].Id {
        return fmt.Errorf("followings %d and %d had the same ids", i + 1, j + 1)
      }
    }
  }
  for i := 0; i < len(user.ChatIds); i++ {
    for j := i + 1; j < len(user.ChatIds); j++ {
      if user.ChatIds[j] == user.ChatIds[i] {
        return fmt.Errorf("user has duplicate chat_id %d", user.ChatIds[i])
      }
    }
  }
  user.Pretty()
  return nil
}

func (user *User) Markdown() string {
  return "[" + tools.Escape(user.Name) + "](https://twitter.com/" + tools.EscapeLink(user.Username) + ")"
}

func (user *User) Copy() *User {
  var newUser User
  newUser.Id = user.Id
  newUser.Username = user.Username
  newUser.Name = user.Name
  newUser.Followings = append(newUser.Followings, user.Followings...)
  newUser.ChatIds = append(newUser.ChatIds, user.ChatIds...)
  return &newUser
}
