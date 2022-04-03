package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "io"
    "log"
    "net/http"
    "net/url"
    "os"
    "time"
)

type User struct {
    Id string
    Username string
    Name string
}

type UserData struct {
    User
    Following []User
}

var (
    TWITTER_TOKEN string
    TWITTER_API string
    TELEGRAM_API string
    TELEGRAM_CHAT_IDS []string
)

func readTelegramChatIds() ([]string, error) {
    var telegramChatIds []string
    file, err := os.Open("telegram_chat_ids")
    if err != nil {
        return telegramChatIds, err
    }
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        telegramChatIds = append(telegramChatIds, scanner.Text())
    }
    file.Close()
    return telegramChatIds, scanner.Err()
}

func readData() ([]UserData, error) {
    var data []UserData
    file, err := os.Open("data.json")
    if err != nil {
        return data, err
    }
    bytes, err := io.ReadAll(file)
    file.Close()
    if err != nil {
        return data, err
    }
    err = json.Unmarshal(bytes, &data)
    return data, err
}

func readUsernames(usernamesFile string) ([]string, error) {
    var usernames []string
    file, err := os.Open(usernamesFile)
    if err != nil {
        return usernames, err
    }
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        usernames = append(usernames, scanner.Text())
    }
    file.Close()
    return usernames, scanner.Err()
}

func getUsers(usernames []string) ([]User, error) {
    var users []User
    client := &http.Client{}
    for i := 0; i < len(usernames); i += 100 {
        query := TWITTER_API + "by?usernames=" + usernames[i]
        for j := 1; j < 100 && i + j < len(usernames); j++ {
            query += "," + usernames[i + j]
        }
        req, err := http.NewRequest("GET", query, nil)
        if err != nil {
            return users, err
        }
        req.Header.Add("Authorization", "Bearer " + TWITTER_TOKEN)
        resp, err := client.Do(req)
        if err != nil {
            return users, err
        }
        body, err := io.ReadAll(resp.Body)
        resp.Body.Close()
        if err != nil {
            return users, err
        }
        var jsonBody struct {
            Data []User
        }
        err = json.Unmarshal(body, &jsonBody)
        if err != nil {
            return users, err
        }
        users = append(users, jsonBody.Data...)
    }
    return users, nil
}

func getFollowingUsers(userId string) ([]User, error) {
    var users []User
    client := &http.Client{}
    req, err := http.NewRequest("GET", TWITTER_API + userId + "/following?max_results=1000", nil)
    if err != nil {
        return users, err
    }
    req.Header.Add("Authorization", "Bearer " + TWITTER_TOKEN)
    resp, err := client.Do(req)
    if err != nil {
        return users, err
    }
    body, err := io.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
        return users, err
    }
    var jsonBody struct {
        Data []User
    }
    err = json.Unmarshal(body, &jsonBody)
    return jsonBody.Data, err
}

func getData(users []User, delay time.Duration) ([]UserData, error) {
    var data []UserData
    for _, user := range users {
        time.Sleep(delay)
        log.Println(user.Name + " (https://twitter.com/" + user.Username + ")")
        followingUsers, err := getFollowingUsers(user.Id)
        if err != nil {
            return data, err
        }
        data = append(data, UserData{user, followingUsers})
    }
    return data, nil
}

func writeData(data []UserData) error {
    bytes, err := json.Marshal(data)
    if err != nil {
        return err
    }
    file, err := os.Create("data.json")
    if err != nil {
        return err
    }
    _, err = file.Write(bytes)
    file.Close()
    return err
}

func prepare() ([]UserData, error) {
    log.Println("preparing")

    TWITTER_TOKEN = os.Getenv("TWITTER_TOKEN")
    TWITTER_API = "https://api.twitter.com/2/users/"
    TELEGRAM_API = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_TOKEN") + "/"

    var data []UserData

    var err error
    TELEGRAM_CHAT_IDS, err = readTelegramChatIds()
    if err != nil {
        return data, err
    }

    var usernamesFile string
    flag.StringVar(&usernamesFile, "f", "", "path to usernames file if you want to reload all usernames")
    flag.Parse()

    if usernamesFile == "" {
        log.Println("usernames file not specified, trying to load saved data from \"data.json\"")
        return readData()
    }

    log.Println("usernames file \"" + usernamesFile + "\", trying to load usernames")

    usernames, err := readUsernames(usernamesFile)
    if err != nil {
        return data, err
    }

    log.Println("getting users by usernames")

    users, err := getUsers(usernames)
    if err != nil {
        return data, err
    }

    log.Println("creating new data")

    data, err = getData(users, 60 * time.Second)
    if err != nil {
        return data, err
    }

    return data, writeData(data)
}

func diff(usersA, usersB []User) []User {
    var users []User
    for _, userA := range usersA {
        found := false
        for _, userB := range usersB {
            if userA.Id == userB.Id {
                found = true
                break
            }
        }
        if !found {
            users = append(users, userA)
        }
    }
    return users
}

func sendMessage(chatId, text string) error {
    resp, err := http.PostForm(
        TELEGRAM_API + "sendMessage",
        url.Values{
            "chat_id": {chatId},
            "text": {text},
            "disable_web_page_preview": {"true"},
        },
    )
    if err != nil {
        return err
    }
    resp.Body.Close()
    return nil
}

func sendUpdates(updateUser *UserData, text string, users []User) error {
    text = updateUser.Name + " (https://twitter.com/" + updateUser.Username + ") " + text
    for _, user := range users {
        text += "\n" + user.Name + " (https://twitter.com/" + user.Username + ")"
    }
    log.Println(text)
    for _, chatId := range TELEGRAM_CHAT_IDS {
        err := sendMessage(chatId, text)
        if err != nil {
            return err
        }
    }
    return nil
}

func update(user *UserData, followingUsers []User) error {
    users := diff(user.Following, followingUsers)
    if len(users) > 0 {
        err := sendUpdates(user, "unfollowed:", users)
        if err != nil {
            return err
        }
    }
    users = diff(followingUsers, user.Following)
    if len(users) > 0 {
        err := sendUpdates(user, "followed:", users)
        if err != nil {
            return err
        }
    }
    user.Following = followingUsers
    return nil
}

func monitor(data []UserData, delay time.Duration) error {
    log.Println("monitoring")
    userCount := len(data)
    for i := 0; i < userCount; i = (i + 1) % userCount {
        time.Sleep(delay)
        followingUsers, err := getFollowingUsers(data[i].Id)
        if err != nil {
            return err
        }
        err = update(&data[i], followingUsers)
        if err != nil {
            return err
        }
        err = writeData(data)
        if err != nil {
            return err
        }
    }
    return nil
}

func main() {
    data, err := prepare()
    if err != nil {
        panic(err)
    }
    log.Printf("working with users:")
    for _, user := range data {
        log.Println(user.Name + " (https://twitter.com/" + user.Username + ")")
    }
    err = monitor(data, 60 * time.Second)
    if err != nil {
        panic(err)
    }
}
