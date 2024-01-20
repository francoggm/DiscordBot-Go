package bot

import (
	"discord-bot/schedule"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func getCustomInteractionResponse(content string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}
}

func getCustomInteractionError(err error) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: cases.Title(language.English, cases.NoLower).String(err.Error()) + "!",
		},
	}
}

var helloWorldHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, getCustomInteractionResponse("Hello there, i'm a bot!"))
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
		s.InteractionRespond(i.Interaction, getCustomInteractionError(err))
	} else {
		res := fmt.Sprintf("Successfuly schedule \"%s\" for %s %s", sc.Appointment, sc.Day, sc.Hours)
		s.InteractionRespond(i.Interaction, getCustomInteractionResponse(res))
	}
}
