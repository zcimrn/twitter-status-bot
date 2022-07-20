package data

import (
	"sort"
)

func (data *Data) GetChats() []Chat {
	data.mutex.RLock()
	chats := append([]Chat{}, data.Chats...)
	data.mutex.RUnlock()
	return chats
}

func (data *Data) HasChat(id int) bool {
	data.mutex.RLock()
	i := sort.Search(len(data.Chats), func(i int) bool {
		return data.Chats[i].Id >= id
	})
	found := i < len(data.Chats) && data.Chats[i].Id == id
	data.mutex.RUnlock()
	return found
}

func (data *Data) AddChat(chat *Chat) bool {
	data.mutex.Lock()
	i := sort.Search(len(data.Chats), func(i int) bool {
		return data.Chats[i].Id >= chat.Id
	})
	if i < len(data.Chats) && data.Chats[i].Id == chat.Id {
		data.Chats[i] = *chat
		data.save()
		data.mutex.Unlock()
		return true
	}
	data.Chats = append(data.Chats[:i], append([]Chat{*chat}, data.Chats[i:]...)...)
	data.save()
	data.mutex.Unlock()
	return false
}

func (data *Data) DeleteChat(id int) bool {
	data.mutex.Lock()
	i := sort.Search(len(data.Chats), func(i int) bool {
		return data.Chats[i].Id >= id
	})
	if i == len(data.Chats) || data.Chats[i].Id != id {
		data.mutex.Unlock()
		return false
	}
	data.Chats = append(data.Chats[:i], data.Chats[i+1:]...)
	data.save()
	data.mutex.Unlock()
	return true
}
