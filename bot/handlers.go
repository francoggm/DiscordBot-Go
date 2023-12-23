package bot

import (
	"discord-bot/weather"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func receiveMessage(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == session.State.User.ID {
		return
	}

	switch {
	case strings.Contains(msg.Content, "weather"):
		session.ChannelMessageSend(msg.ChannelID, "I can help you with that! Use '!zip<zip code>'")
	case strings.Contains(msg.Content, "bot"):
		session.ChannelMessageSend(msg.ChannelID, "Hello world!")
	case strings.Contains(msg.Content, "!zip"):
		currentWeather := weather.GetWeather(msg.Content)
		session.ChannelMessageSendComplex(msg.ChannelID, currentWeather)
	}
}
