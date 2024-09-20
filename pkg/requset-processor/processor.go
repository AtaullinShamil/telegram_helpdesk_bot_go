package requst_processor

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RequestProcessor struct {
	Bot *tgbotapi.BotAPI
	Db  map[int64]Request
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

	Processor.Db = make(map[int64]Request, 0)

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
