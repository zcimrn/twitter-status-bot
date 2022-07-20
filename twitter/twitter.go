package twitter

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

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
	for nextToken := ""; ; {
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
				ResultCount int    `json:"result_count"`
				NextToken   string `json:"next_token"`
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
	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})
	return users, nil
}
