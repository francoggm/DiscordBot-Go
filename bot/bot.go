package bot

import (
	"discord-bot/config"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	token              string
	appID              string
	guildID            string
	session            *discordgo.Session
	registeredCommands []*discordgo.ApplicationCommand
}

func New() (*Bot, error) {
	cfg := config.GetConfig()

	s, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, err
	}

	return &Bot{
		token:              cfg.BotToken,
		appID:              cfg.AppID,
		guildID:            cfg.GuildID,
		session:            s,
		registeredCommands: make([]*discordgo.ApplicationCommand, len(commands)),
	}, nil
}

func (b *Bot) Run() {
	b.session.AddHandler(receiveLogin)
	b.session.AddHandler(receiveInteraction)

	err := b.session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session : Error=%v", err)
	}
	defer b.session.Close()

	for i, command := range commands {
		cmd, err := b.session.ApplicationCommandCreate(b.appID, b.guildID, command)
		if err != nil {
			log.Panicf("Failed to add command=%s : Error=%s", command.Name, err.Error())
		}

		b.registeredCommands[i] = cmd
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	for _, command := range b.registeredCommands {
		err := b.session.ApplicationCommandDelete(b.appID, b.guildID, command.ID)
		if err != nil {
			log.Panicf("Failed to delete command=%s : Error=%s", command.Name, err.Error())
		}
	}
}
