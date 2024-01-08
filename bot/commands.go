package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "hello-world",
			Description: "Returns a hello string",
		},
		{
			Name:        "schedule",
			Description: "Schedule some appointment",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "appointment",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "Description of the commitment",
					Required:    true,
				},
				{
					Name:        "day",
					Description: "The day the appointment will be (dd/mm/yy)",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "hour",
					Description: "The time the appointment will be (hh:mm AM/PM)",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hello-world": helloWorld,
		"schedule":    schedule,
	}
)

var helloWorld = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello there, i'm a bot",
		},
	})
}

var schedule = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var (
		appointment string
		day         string
		hours       string
	)

	options := i.ApplicationCommandData().Options
	for _, option := range options {
		switch option.Name {
		case "appointment":
			appointment = option.StringValue()
		case "day":
			day = option.StringValue()
		case "hour":
			hours = option.StringValue()
		}
	}

	res := fmt.Sprintf("Appointment=%s : Day=%s : Hour=%s", appointment, day, hours)

	log.Println(res)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: res,
		},
	})
}
