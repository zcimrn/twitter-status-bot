package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func TestToken(token string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/by/username/zcimrn", nil)
	if err != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

func api(query string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+getToken())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("%s - %s", query, resp.Status)
	}
	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var jsonResp struct {
		Errors []struct {
			Detail string `json:"detail"`
		} `json:"errors"`
	}
	err = json.Unmarshal(respBody, &jsonResp)
	if err != nil {
		return nil, err
	}
	if len(jsonResp.Errors) > 0 {
		return nil, fmt.Errorf(jsonResp.Errors[0].Detail)
	}
	return respBody, nil
}
