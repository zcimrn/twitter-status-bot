package telegram

import (
  "bytes"
  "encoding/json"
  "fmt"
  "io"
  "net/http"

  "github.com/zcimrn/twitter-status-bot/config"
)

var (
  Config *config.Config
)

func TestToken(token string) bool {
  resp, err := http.Get("https://api.telegram.org/bot" + token + "/getMe")
  if err != nil {
    return false
  }
  resp.Body.Close()
  if resp.StatusCode != 200 {
    return false
  }
  return true
}

func api(method string, reqBody []byte) ([]byte, error) {
  url := "https://api.telegram.org/bot" + Config.GetTelegramToken() + "/" + method
  resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
  if err != nil {
    return nil, err
  }
  respBody, err := io.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    return nil, err
  }
  var jsonResp struct {
    Ok bool `json:"ok"`
    Desc string `json:"description"`
  }
  err = json.Unmarshal(respBody, &jsonResp)
  if err != nil {
    return nil, err
  }
  if !jsonResp.Ok {
    return nil, fmt.Errorf(jsonResp.Desc)
  }
  return respBody, nil
}