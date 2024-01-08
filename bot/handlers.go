package bot

import (
	"github.com/bwmarrin/discordgo"
)

func receiveInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		handler(s, i)
	}
}
