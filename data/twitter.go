package data

import (
	"github.com/zcimrn/twitter-status-bot/twitter"
)

func (data *Data) GetTwitterToken() string {
	data.mutex.RLock()
	twitterToken := data.TwitterToken
	data.mutex.RUnlock()
	return twitterToken
}

func (data *Data) SetTwitterToken(twitterToken string) {
	data.mutex.Lock()
	data.TwitterToken = twitterToken
	twitter.SetToken(twitterToken)
	data.save()
	data.mutex.Unlock()
}
