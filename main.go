package main

import (
	"discord-bot/bot"
	"discord-bot/config"
	"log"
)

func main() {
	err := config.Read()
	if err != nil {
		log.Fatalf("Error reading configs, error=%s", err.Error())
	}

	cfg := config.GetConfig()

	bot, err := bot.Config(cfg.BotToken, cfg.OpenWeatherToken)
	if err != nil {
		log.Fatalf("Error setting bot, error=%s", err.Error())
	}

	bot.Run()
}
