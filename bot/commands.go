package bot

import (
	"discord-bot/schedule"
	"fmt"

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

func customInteractionResponse(content string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}
}

func customInteractionError(err error) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: err.Error(),
		},
	}
}

var helloWorldHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello there, i'm a bot",
		},
	})
}

var scheduleHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var sc schedule.Schedule

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
		s.InteractionRespond(i.Interaction, customInteractionError(err))
	} else {
		res := fmt.Sprintf("Successfuly schedule \"%s\" for %s %s", sc.Appointment, sc.Day, sc.Hours)
		s.InteractionRespond(i.Interaction, customInteractionResponse(res))
	}
}
