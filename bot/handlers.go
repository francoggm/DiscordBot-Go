package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func receiveLogin(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as %v#%v", s.State.User.Username, s.State.User.Discriminator)
}

func receiveInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		handler(s, i)
	}
}
