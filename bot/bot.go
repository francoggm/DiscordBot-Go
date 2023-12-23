package bot

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	botToken         string
	openWeatherToken string
	session          *discordgo.Session
}

func Config(botToken string, openWeatherToken string) (*Bot, error) {
	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	return &Bot{
		botToken:         botToken,
		openWeatherToken: openWeatherToken,
		session:          s,
	}, nil
}

func (b *Bot) Run() {
	b.session.AddHandler(receiveMessage)

	b.session.Open()
	defer b.session.Close()

	fmt.Println("Running")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
