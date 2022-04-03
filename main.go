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
    log.Printf("reading telegram chat ids from \"telegram_chat_ids\"")
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
    log.Printf("reading data from \"data.json\"")
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
    log.Printf("reading usernames from \"%s\" file", usernamesFile)
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

func pretty(name, username string) string {
    return "[" + name + "](https://twitter.com/" + username + ")"
}

func getUsers(usernames []string) ([]User, error) {
    log.Printf("getting users by usernames")
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
        for _, user := range jsonBody.Data {
            log.Printf(pretty(user.Name, user.Username))
            users = append(users, user)
        }
    }
    return users, nil
}

func getFollowingUsers(userId string, delay time.Duration) ([]User, error) {
    var users []User
    client := &http.Client{}
    for nextToken := "";; {
        time.Sleep(delay)
        query := TWITTER_API + userId + "/following?max_results=1000"
        if nextToken != "" {
            query += "&pagination_token=" + nextToken
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
            Data []User `json:"data"`
            Meta struct {
                ResultCount int `json:"result_count"`
                NextToken string `json:"next_token"`
            } `json:"meta"`
        }
        err = json.Unmarshal(body, &jsonBody)
        if err != nil {
            return users, err
        }
        users = append(users, jsonBody.Data...)
        nextToken = jsonBody.Meta.NextToken
        if nextToken == "" {
            break
        }
    }
    return users, nil
}

func getData(users []User, delay time.Duration) ([]UserData, error) {
    log.Printf("getting users data")
    var data []UserData
    for _, user := range users {
        log.Printf("%s getting followings", pretty(user.Name, user.Username))
        followingUsers, err := getFollowingUsers(user.Id, delay)
        if err != nil {
            return data, err
        }
        log.Printf("%s got %d followings", pretty(user.Name, user.Username), len(followingUsers))
        data = append(data, UserData{user, followingUsers})
    }
    return data, nil
}

func writeData(data []UserData) error {
    log.Printf("writing data to \"data.json\"")
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
    log.Printf("preparing")

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
        log.Printf("usernames file not specified")
        return readData()
    }

    usernames, err := readUsernames(usernamesFile)
    if err != nil {
        return data, err
    }

    users, err := getUsers(usernames)
    if err != nil {
        return data, err
    }

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
            "parse_mode": {"MarkdownV2"},
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
    log.Printf("sending updates")
    text = pretty(updateUser.Name, updateUser.Username) + " " + text + ":"
    for _, user := range users {
        text += "\n" + pretty(user.Name, user.Username)
    }
    log.Printf(text)
    for _, chatId := range TELEGRAM_CHAT_IDS {
        err := sendMessage(chatId, text)
        if err != nil {
            return err
        }
    }
    return nil
}

func update(user *UserData, followingUsers []User) error {
    log.Printf("%s updating", pretty(user.Name, user.Username))
    /*
    users := diff(user.Following, followingUsers)
    if len(users) > 0 {
        err := sendUpdates(user, "unfollowed", users)
        if err != nil {
            return err
        }
    }
    */
    users := diff(followingUsers, user.Following)
    if len(users) > 0 {
        err := sendUpdates(user, "followed", users)
        if err != nil {
            return err
        }
    }
    user.Following = followingUsers
    return nil
}

func monitor(data []UserData, delay time.Duration) error {
    log.Printf("monitoring")
    userCount := len(data)
    for i := 0; i < userCount; i = (i + 1) % userCount {
        log.Printf("%s getting followings", pretty(data[i].Name, data[i].Username))
        followingUsers, err := getFollowingUsers(data[i].Id, delay)
        if err != nil {
            return err
        }
        log.Printf("%s got %d followings", pretty(data[i].Name, data[i].Username), len(followingUsers))
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
        log.Printf(pretty(user.Name, user.Username))
    }
    err = monitor(data, 60 * time.Second)
    if err != nil {
        panic(err)
    }
}
