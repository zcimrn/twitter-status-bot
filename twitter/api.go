package twitter

import (
  "encoding/json"
  "errors"
  "io"
  "net/http"

  "github.com/zcimrn/twitter-status-bot/config"
)

var Config *config.Config

func TestToken(token string) bool {
  client := &http.Client{}
  req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/by/username/zcimrn", nil)
  if err != nil {
    return false
  }
  req.Header.Add("Authorization", "Bearer " + token)
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
  req.Header.Add("Authorization", "Bearer " + Config.GetTwitterToken())
  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  if resp.StatusCode != 200 {
    resp.Body.Close()
    return nil, errors.New(query + " - " + resp.Status)
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
    return nil, errors.New(jsonResp.Errors[0].Detail)
  }
  return respBody, nil
}
