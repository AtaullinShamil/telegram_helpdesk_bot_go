package main

import (
	"log"

	processor "github.com/AtaullinShamil/telegram_helpdesk_bot_go/pkg/requset-processor"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	Proc, err := processor.NewRequestProcessor()
	if err != nil {
		log.Panic(err)
	}
	//Proc.Bot.Debug = true
	log.Printf("Authorized on account %s", Proc.Bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"callback_query", "message"}

	updates := Proc.Bot.GetUpdatesChan(u)

	for update := range updates {
		Proc.HandleUpdate(update)
	}
}
