package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	configor.Load(&Config, "config.yml")
	price, _ := getBitcoinPriceIndex(Config.PriceIndexURI)

	fmt.Printf("Bitcoin Index: %.3f $\n", price.Usd.RateFloat)

	bot, err := tgbotapi.NewBotAPI(Config.TelegramBotToken)
	if err != nil {
		log.Panic("telegram bot: %s", err)
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := update.Message.Text
		log.Printf("[%s] %s", update.Message.From.UserName, message)

		reply := "I don't understand you"

		switch message {
		case "/start":
			reply = "Type /price_index"
		case "/price_index":
			reply = fmt.Sprintf("Bitcoin Index: %.3f $\n", price.Usd.RateFloat)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		bot.Send(msg)
	}
}
