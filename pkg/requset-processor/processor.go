package requst_processor

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RequestProcessor struct {
	Bot *tgbotapi.BotAPI

	admins    map[string]int64
	passwords map[string]string

	Db            map[int64]Request
	OpenTicketsDb map[int64]Request
}

func NewRequestProcessor() (*RequestProcessor, error) {
	Processor := &RequestProcessor{}

	token, exists := os.LookupEnv("BOTTOKEN")
	if !exists {
		return nil, fmt.Errorf("there isn't bot token env")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	Processor.Bot = bot

	Processor.admins = map[string]int64{
		"Support": 0,
		"IT":      0,
		"Billing": 0,
	}

	supportPassword, exists := os.LookupEnv("SUPPORTPASSWORD")
	if !exists {
		return nil, fmt.Errorf("there isn't bot token env")
	}

	itPassword, exists := os.LookupEnv("ITPASSWORD")
	if !exists {
		return nil, fmt.Errorf("there isn't bot token env")
	}

	billingPassword, exists := os.LookupEnv("BILLINGPASSWORD")
	if !exists {
		return nil, fmt.Errorf("there isn't bot token env")
	}

	Processor.passwords = map[string]string{
		"Support": supportPassword,
		"IT":      itPassword,
		"Billing": billingPassword,
	}

	Processor.Db = make(map[int64]Request, 0)
	Processor.OpenTicketsDb = make(map[int64]Request, 0)

	return Processor, nil
}

func (p *RequestProcessor) HandleUpdate(update tgbotapi.Update) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			log.Printf("recovered from panic: %v", panicValue)
		}
	}()

	if update.CallbackQuery != nil {
		p.handleCallback(update.CallbackQuery)
		return
	}

	if update.Message == nil {
		return
	}

	switch update.Message.Command() {
	case "start":
		p.handleStart(update.Message)
		return
	case "new":
		p.handleNew(update.Message)
		return
	}

	p.handleMessages(update.Message)
}
