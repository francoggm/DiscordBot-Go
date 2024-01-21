package bot

import (
	"discord-bot/schedule"

	"github.com/bwmarrin/discordgo"
)

// Interactions

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
					Description: "Description of the appointment",
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
					Description: "The time the appointment will be (hh:mm)",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hello-world": helloWorldHandler,
		"schedule":    scheduleHandler,
	}
)

// Handlers functions

var helloWorldHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello there! I'm a bot",
		},
	})
}

var scheduleHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var sc schedule.Schedule

	// user is nil if interaction is from server
	if i.User != nil {
		sc.UserID = i.User.ID
	} else {
		sc.UserID = i.Member.User.ID
	}

	options := i.ApplicationCommandData().Options
	for _, option := range options {
		switch option.Name {
		case "appointment":
			sc.Appointment = option.StringValue()
		case "day":
			sc.Day = option.StringValue()
		case "hour":
			sc.Hours = option.StringValue()
		}
	}

	err := schedule.ScheduleAppointment(sc)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Error in schedule!",
						Description: err.Error(),
						Color:       0xB02506,
					},
				},
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Successfuly schedule!",
						Description: "\"" + sc.Appointment + "\" will be alerted in " + sc.Day + " " + sc.Hours,
						Color:       0x0EB625,
					},
				},
			},
		})
	}
}
