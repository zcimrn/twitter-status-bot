package data

import (
	"sort"
)

func (data *Data) GetAdmins() []Admin {
	data.mutex.RLock()
	admins := append([]Admin{}, data.Admins...)
	data.mutex.RUnlock()
	return admins
}

func (data *Data) HasAdmin(id int) bool {
	data.mutex.RLock()
	i := sort.Search(len(data.Admins), func(i int) bool {
		return data.Admins[i].Id >= id
	})
	found := i < len(data.Admins) && data.Admins[i].Id == id
	data.mutex.RUnlock()
	return found
}

func (data *Data) AddAdmin(admin *Admin) bool {
	data.mutex.Lock()
	i := sort.Search(len(data.Admins), func(i int) bool {
		return data.Admins[i].Id >= admin.Id
	})
	if i < len(data.Admins) && data.Admins[i].Id == admin.Id {
		data.Admins[i] = *admin
		data.save()
		data.mutex.Unlock()
		return true
	}
	data.Admins = append(data.Admins[:i], append([]Admin{*admin}, data.Admins[i:]...)...)
	data.save()
	data.mutex.Unlock()
	return false
}

func (data *Data) DeleteAdmin(id int) bool {
	data.mutex.Lock()
	i := sort.Search(len(data.Admins), func(i int) bool {
		return data.Admins[i].Id >= id
	})
	if i == len(data.Admins) || data.Admins[i].Id != id {
		data.mutex.Unlock()
		return false
	}
	data.Admins = append(data.Admins[:i], data.Admins[i+1:]...)
	data.save()
	data.mutex.Unlock()
	return true
}
