package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	go startPriceSync(config.Bitcoin.Uri, time.Duration(config.Bitcoin.Timeout))

	bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.Fatalf("telegram bot: %s", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := update.Message.Text

		currentUser := getUser(update.Message.From.ID)
		if currentUser == nil {
			setUser(update.Message.From.ID, update.Message.From.UserName)
		}

		reply := "Unknown command. Type /help"
		commands := strings.Fields(message)

		switch commands[0] {
		case "/start":
			reply = fmt.Sprintf("Type /help for help. Bitcoin indexes are updating every %d seconds.", config.Bitcoin.Timeout)

		case "/help":
			reply = fmt.Sprint("/prices - show price indexes\n/subscribe USD|EUR|GPB - make a subscription to USD|EUR|GPB \n/unsubscribe USD|EUR|GPB - unsubscribe from USD|EUR|GPB")

		case "/subscribe":
			if subscribe(update.Message.From.ID, commands[1]) {
				reply = fmt.Sprintf("Subscribed to %s\n", commands[1])
			} else {
				reply = fmt.Sprintf("Unknown currency code: %s", commands[1])
			}

		case "/unsubscribe":
			if unsubscribe(update.Message.From.ID, commands[1]) {
				reply = fmt.Sprintf("Unsubscribed from %s", commands[1])
			} else {
				reply = fmt.Sprintf("Unknown currency code: %s", commands[1])
			}

		case "/prices":
			subsriptions := getSubscriptionsByUserID(update.Message.From.ID)
			replyBuf := new(bytes.Buffer)
			for _, element := range subsriptions {
				replyBuf.WriteString(
					fmt.Sprintf("Bitcoin index: %.3f %s\n", element.Currency.Rate, element.Currency.Code))
			}
			reply = replyBuf.String()

		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
