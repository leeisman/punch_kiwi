package main

import (
	"fmt"
	"log"
	"punch_kiwi/chromedp"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	Token string = "2005698212:AAFBe4GtcdqIfiSCSqmi5wYgBCv7IKDyJuw"
	URL   string = "https://cloud.nueip.com/login/83663709"
)

var (
	indexUsername map[string]*User
)

func main() {
	indexUsername = make(map[string]*User, 0)
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	fmt.Println("config:", config)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		command := update.Message.Text
		username := update.Message.From.UserName
		ok, user := checkEnableUser(config, username)
		rsp := "received"
		if !ok {
			rsp = fmt.Sprintf("username: %s disabled\n", username)
		} else {
			err := commandExecute(command, user)
			if err != nil {
				rsp = err.Error()
			}
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, rsp)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func checkEnableUser(config *Config, username string) (bool, *User) {
	if len(indexUsername) < 1 {
		for _, user := range config.Users {
			indexUsername[user.Username] = user
		}
	}
	user, ok := indexUsername[username]
	return ok, user
}

func commandExecute(command string, user *User) error {
	go func() {
		chromedp.NewHooker(URL, user.ID, user.Password, command)
	}()
	fmt.Println(user.Username, command)
	return nil
}
