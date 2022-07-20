package twitter

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zcimrn/twitter-status-bot/tools"
)

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Initialized bool   `json:"initialized"`
	Followings  []User `json:"followings"`
	ChatIds     []int  `json:"chat_ids"`
}

func (user *User) Validate() error {
	if user.Id == "" {
		return fmt.Errorf("user has empty id")
	}
	if _, err := strconv.Atoi(user.Id); err != nil {
		return fmt.Errorf("user has non int id")
	}
	user.Username = strings.ToLower(user.Username)
	for i := 0; i < len(user.Followings); i++ {
		if err := user.Followings[i].Validate(); err != nil {
			return fmt.Errorf("following %d error: '%s'", err)
		}
		for j := i + 1; j < len(user.Followings); j++ {
			if user.Followings[j].Id == user.Followings[i].Id {
				return fmt.Errorf("followings %d and %d have the same ids", i+1, i+1)
			}
		}
	}
	sort.Slice(user.Followings, func(i, j int) bool {
		return user.Followings[i].Id < user.Followings[j].Id
	})
	for i := 0; i < len(user.ChatIds); i++ {
		for j := i + 1; j < len(user.ChatIds); j++ {
			if user.ChatIds[j] == user.ChatIds[i] {
				return fmt.Errorf("duplicate chat_id %d", user.ChatIds[i])
			}
		}
	}
	sort.Ints(user.ChatIds)
	return nil
}

func (user *User) Markdown() string {
	return fmt.Sprintf("[%s](https://twitter.com/%s)", tools.Escape(user.Name), tools.EscapeLink(user.Username))
}

func (user *User) HasFollowing(id string) bool {
	i := sort.Search(len(user.Followings), func(i int) bool {
		return user.Followings[i].Id >= id
	})
	return i < len(user.Followings) && user.Followings[i].Id == id
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
			if !user.HasFollowing(newUser.Followings[i].Id) {
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
	i := sort.Search(len(user.ChatIds), func(i int) bool {
		return user.ChatIds[i] >= id
	})
	if i == len(user.ChatIds) || user.ChatIds[i] != id {
		return false
	}
	user.ChatIds = append(user.ChatIds[:i], user.ChatIds[i+1:]...)
	return true
}
