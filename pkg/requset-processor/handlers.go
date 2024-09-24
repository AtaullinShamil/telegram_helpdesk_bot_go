package requst_processor

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Please, choose department")
	inlineKey := &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{tgbotapi.NewInlineKeyboardButtonData("Support", "support:")},
			{tgbotapi.NewInlineKeyboardButtonData("IT", "it:")},
			{tgbotapi.NewInlineKeyboardButtonData("Billing", "billing:")},
		},
	}
	msg.ReplyMarkup = inlineKey
	_, err := p.Bot.Send(msg)
	if err != nil {
		log.Printf("Error send message: %v", err)
		return
	}

	id := inputMessage.From.ID
	p.Db[id] = Request{UserId: id, ChatId: inputMessage.Chat.ID}
}

func (p *RequestProcessor) handleCallback(callback *tgbotapi.CallbackQuery) {
	splited := strings.Split(callback.Data, ":")
	action := splited[0]

	switch action {
	case "support":
		id := callback.From.ID
		val := p.Db[id]
		val.Department = "Support"
		val.Status.IsDepartment = true
		p.Db[id] = val

		chatID := callback.Message.Chat.ID

		deleted := tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)
		response, err := p.Bot.Request(deleted)
		if err != nil {
			log.Printf("Error sending deleted message: %v", err)
			return
		}
		if !response.Ok {
			log.Printf("Error delete message: %v", err)
			return
		}

		msg := tgbotapi.NewMessage(chatID, "You choose Support department")
		_, err = p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}

		msgTittle := tgbotapi.NewMessage(chatID, "Please, write Tittle of ticket")
		_, err = p.Bot.Send(msgTittle)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}

	case "it":
		id := callback.From.ID
		val := p.Db[id]
		val.Department = "IT"
		val.Status.IsDepartment = true
		p.Db[id] = val

		chatID := callback.Message.Chat.ID

		deleted := tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)
		response, err := p.Bot.Request(deleted)
		if err != nil {
			log.Printf("Error sending deleted message: %v", err)
			return
		}
		if !response.Ok {
			log.Printf("Error delete message: %v", err)
			return
		}

		msg := tgbotapi.NewMessage(chatID, "You choose IT department")
		_, err = p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}

		msgTittle := tgbotapi.NewMessage(chatID, "Please, write Tittle of ticket")
		_, err = p.Bot.Send(msgTittle)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
	case "billing":
		id := callback.From.ID
		val := p.Db[id]
		val.Department = "Billing"
		val.Status.IsDepartment = true
		p.Db[id] = val

		chatID := callback.Message.Chat.ID

		deleted := tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)
		response, err := p.Bot.Request(deleted)
		if err != nil {
			log.Printf("Error sending deleted message: %v", err)
			return
		}
		if !response.Ok {
			log.Printf("Error delete message: %v", err)
			return
		}

		msg := tgbotapi.NewMessage(chatID, "You choose Support department")
		_, err = p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}

		msgTittle := tgbotapi.NewMessage(chatID, "Please, write Tittle of ticket")
		_, err = p.Bot.Send(msgTittle)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
	case "accept":
		chatID := callback.Message.Chat.ID
		deleted := tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)
		response, err := p.Bot.Request(deleted)
		if err != nil {
			log.Printf("Error sending deleted message: %v", err)
			return
		}
		if !response.Ok {
			log.Printf("Error delete message: %v", err)
			return
		}

		id := callback.From.ID
		val, ok := p.Db[id]
		if !ok {
			msg := tgbotapi.NewMessage(chatID, "There isn't ticket")
			_, err := p.Bot.Send(msg)
			if err != nil {
				log.Printf("Error send message: %v", err)
				return
			}
			return
		}
		msg := tgbotapi.NewMessage(chatID, "Ticket assepted")
		_, err = p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}

		delete(p.Db, id)
		p.OpenTicketsDb[id] = val //add postgres

		adminMsg := tgbotapi.NewMessage(p.admins[val.Department], fmt.Sprintf("UserID : %d\nDepartment : %s\nTitle : %s\nDiscription: %s", val.UserId, val.Department, val.Tittle, val.Discription))
		inlineKey := &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{tgbotapi.NewInlineKeyboardButtonData("Answer the ticket", fmt.Sprintf("answer:%d", val.ChatId))},
			},
		}
		adminMsg.ReplyMarkup = inlineKey
		_, err = p.Bot.Send(adminMsg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
	case "delete":
		chatID := callback.Message.Chat.ID
		deleted := tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID)
		response, err := p.Bot.Request(deleted)
		if err != nil {
			log.Printf("Error sending deleted message: %v", err)
			return
		}
		if !response.Ok {
			log.Printf("Error delete message: %v", err)
			return
		}

		id := callback.From.ID
		_, ok := p.Db[id]
		if !ok {
			msg := tgbotapi.NewMessage(chatID, "There isn't ticket")
			_, err := p.Bot.Send(msg)
			if err != nil {
				log.Printf("Error send message: %v", err)
				return
			}
			return
		}
		delete(p.Db, id)
		msg := tgbotapi.NewMessage(chatID, "Ticket deleted")
		_, err = p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
	case "answer":
		deleted := tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		response, err := p.Bot.Request(deleted)
		if err != nil {
			log.Printf("Error sending deleted message: %v", err)
			return
		}
		if !response.Ok {
			log.Printf("Error delete message: %v", err)
			return
		}

		chatId, err := strconv.Atoi(splited[1])
		if err != nil {
			log.Printf("Error atoi: %v", err)
			return
		}

		msg := tgbotapi.NewMessage(int64(chatId), "Hello world")
		_, err = p.Bot.Send(msg)
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
	if p.checkPassword(inputMessage.Text, id) {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Admin accepted")
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}
		return
	}
	val, ok := p.Db[id]
	if !ok {
		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Please, use command /new for create a ticket")
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

		msg := tgbotapi.NewMessage(inputMessage.Chat.ID, fmt.Sprintf("Your ticket created :\n\nDepartment : %s\nTittle : %s\nDiscription : %s\n", val.Department, val.Tittle, val.Discription))
		_, err := p.Bot.Send(msg)
		if err != nil {
			log.Printf("Error send message: %v", err)
			return
		}

		inlineKey := &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
				{tgbotapi.NewInlineKeyboardButtonData("Accept", "accept:")},
				{tgbotapi.NewInlineKeyboardButtonData("Delete", "delete:")},
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

func (p *RequestProcessor) checkPassword(input string, id int64) bool {
	for k, v := range p.passwords {
		if input == v {
			p.admins[k] = id
			return true
		}
	}
	return false
}
