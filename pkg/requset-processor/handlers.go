package requst_processor

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startMessage = "Welcome! I am your helpdesk bot!\nYou can use commands :\n/new - for create new tiket"
)

func (p *RequestProcessor) handleStart(inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, startMessage)
	p.Bot.Send(msg)
}

func (p *RequestProcessor) handleNew(inputMessage *tgbotapi.Message) {
	inlineKey := &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{tgbotapi.NewInlineKeyboardButtonData("Support", "1")},
			{tgbotapi.NewInlineKeyboardButtonData("IT", "2")},
			{tgbotapi.NewInlineKeyboardButtonData("Billing", "3")},
		},
	}
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Choose department")
	msg.ReplyMarkup = inlineKey
	_, err := p.Bot.Send(msg)
	if err != nil {
		log.Printf("Error send message: %v", err)
		return
	}

	id := inputMessage.From.ID
	_, ok := p.Db[id]
	if !ok {
		p.Db[id] = Request{}
	}
}

func (p *RequestProcessor) handleCallback(callback *tgbotapi.CallbackQuery) {
	switch callback.Data {
	case "1":
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "You choose Support department")
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
		id := callback.From.ID
		val := p.Db[id]
		val.Department = "Support"
		val.Status.IsDepartment = true
		p.Db[id] = val
		msgTittle := tgbotapi.NewMessage(callback.Message.Chat.ID, "Please, write Tittle of ticket")
		_, err = p.Bot.Send(msgTittle)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
	case "2":
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "You choose IT department")
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
		id := callback.From.ID
		val := p.Db[id]
		val.Department = "IT"
		val.Status.IsDepartment = true
		p.Db[id] = val
		msgTittle := tgbotapi.NewMessage(callback.Message.Chat.ID, "Please, write Tittle of ticket")
		_, err = p.Bot.Send(msgTittle)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
	case "3":
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "You choose Billing department")
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
		id := callback.From.ID
		val := p.Db[id]
		val.Department = "Billing"
		val.Status.IsDepartment = true
		p.Db[id] = val
		msgTittle := tgbotapi.NewMessage(callback.Message.Chat.ID, "Please, write Tittle of ticket")
		_, err = p.Bot.Send(msgTittle)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
	case "accept":
		id := callback.From.ID
		_, ok := p.Db[id]
		if !ok {
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "There isn't ticket")
			_, err := p.Bot.Send(msg)
			if err != nil {
				log.Printf("Error send message: %v", err)
				return
			}
			return
		}
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Ticket assepted")
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
		delete(p.Db, id)
	case "delete":
		id := callback.From.ID
		_, ok := p.Db[id]
		if !ok {
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "There isn't ticket")
			_, err := p.Bot.Send(msg)
			if err != nil {
				log.Printf("Error send message: %v", err)
				return
			}
			return
		}
		delete(p.Db, id)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Ticket deleted")
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
	default:
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Unknown button")
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
	}
}

func (p *RequestProcessor) handleMessages(inputMessage *tgbotapi.Message) {
	id := inputMessage.From.ID
	val, ok := p.Db[id]
	if !ok {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Please, chose command")
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
		return
	}
	if !val.Status.IsTittle {
		val.Tittle = inputMessage.Text
		val.Status.IsTittle = true
		p.Db[id] = val
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Please, write Discription")
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
		return
	}
	if !val.Status.IsDiscription {
		val.Discription = inputMessage.Text
		val.Status.IsDiscription = true
		p.Db[id] = val

		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Your ticket created!\nDepartment : %s\nTittle : %s\nDiscription : %s\n", val.Department, val.Tittle, val.Discription))
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}

		inlineKey := &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{tgbotapi.NewInlineKeyboardButtonData("Accept", "accept")},
				{tgbotapi.NewInlineKeyboardButtonData("Delete", "delete")},
			},
		}
		acceptMsg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Choose action :")
		acceptMsg.ReplyMarkup = inlineKey
		_, err = p.Bot.Send(acceptMsg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
		return
	}
}
