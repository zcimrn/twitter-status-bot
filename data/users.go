package data

import (
	"sort"

	"github.com/zcimrn/twitter-status-bot/twitter"
)

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
			var user twitter.User
			user.Id = data.Users[i].Id
			user.Username = data.Users[i].Username
			user.Name = data.Users[i].Name
			users = append(users, user)
		}
	}
	data.mutex.RUnlock()
	return users
}

func (data *Data) UpdateUser(user *twitter.User) bool {
	data.mutex.Lock()
	i := sort.Search(len(data.Users), func(i int) bool {
		return data.Users[i].Id >= user.Id
	})
	if i == len(data.Users) || data.Users[i].Id != user.Id {
		data.mutex.Unlock()
		return false
	}
	data.Users[i].Username = user.Username
	data.Users[i].Name = user.Name
	data.Users[i].Initialized = user.Initialized
	data.Users[i].Followings = user.Followings
	// hot fix
	var chatIds []int
	for _, id := range data.Users[i].ChatIds {
		i := sort.Search(len(data.Chats), func(i int) bool {
			return data.Chats[i].Id >= id
		})
		if i < len(data.Chats) && data.Chats[i].Id == id {
			chatIds = append(chatIds, id)
		}
	}
	data.Users[i].ChatIds = chatIds
	if len(chatIds) == 0 {
		data.Users = append(data.Users[:i], data.Users[i+1:]...)
		if data.LastIndex >= i {
			data.LastIndex--
		}
	}
	// hot fix
	data.save()
	data.mutex.Unlock()
	return true
}

func (data *Data) AddUser(user *twitter.User) bool {
	data.mutex.Lock()
	i := sort.Search(len(data.Users), func(i int) bool {
		return data.Users[i].Id >= user.Id
	})
	if i < len(data.Users) && data.Users[i].Id == user.Id {
		data.Users[i].Username = user.Username
		data.Users[i].Name = user.Name
		for _, id := range user.ChatIds {
			data.Users[i].AddChatId(id)
		}
		data.save()
		data.mutex.Unlock()
		return true
	}
	data.Users = append(data.Users[:i], append([]twitter.User{*user}, data.Users[i:]...)...)
	if data.LastIndex >= i {
		data.LastIndex++
	}
	data.save()
	data.mutex.Unlock()
	return false
}

func (data *Data) DeleteUser(user *twitter.User) bool {
	data.mutex.Lock()
	i := sort.Search(len(data.Users), func(i int) bool {
		return data.Users[i].Id >= user.Id
	})
	if i == len(data.Users) || data.Users[i].Id != user.Id {
		data.mutex.Unlock()
		return false
	}
	for _, id := range user.ChatIds {
		data.Users[i].DeleteChatId(id)
	}
	if len(data.Users[i].ChatIds) == 0 {
		data.Users = append(data.Users[:i], data.Users[i+1:]...)
		if data.LastIndex >= i {
			data.LastIndex--
		}
	}
	data.save()
	data.mutex.Unlock()
	return true
}

func (data *Data) GetNextUser() *twitter.User {
	data.mutex.Lock()
	if len(data.Users) == 0 {
		data.mutex.Unlock()
		return nil
	}
	data.LastIndex = (data.LastIndex + 1) % len(data.Users)
	// hot fix
	var chatIds []int
	for _, id := range data.Users[data.LastIndex].ChatIds {
		i := sort.Search(len(data.Chats), func(i int) bool {
			return data.Chats[i].Id >= id
		})
		if i < len(data.Chats) && data.Chats[i].Id == id {
			chatIds = append(chatIds, id)
		}
	}
	data.Users[data.LastIndex].ChatIds = chatIds
	// hot fix
	user := data.Users[data.LastIndex]
	data.save()
	data.mutex.Unlock()
	return &user
}
