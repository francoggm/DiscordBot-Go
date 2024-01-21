package schedule

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	fpath = "schedule/schedule.txt"
)

var (
	r *rand.Rand

	colors = []int{
		0x730337,
		0x2160E5,
		0x980CA4,
		0x1FBDC4,
		0xD3A410,
		0x57F1C3,
		0xB034A6,
		0xA40B33,
		0x1C1FC3,
		0xA3C629,
	}

	images = []string{
		"https://static.vecteezy.com/system/resources/previews/022/429/527/non_2x/punctual-being-on-time-or-time-management-work-deadline-or-procrastination-self-discipline-work-efficiency-or-reminder-urgency-or-quick-work-concept-confidence-businessman-jump-over-alarm-clock-vector.jpg",
		"https://img.freepik.com/premium-vector/waste-time-achievement-catch-precious-alarm-clock-management-control-schedule-work_159757-628.jpg",
		"https://img.freepik.com/premium-vector/procrastination-work-time-action-today-timeline-management-postpone-schedule-task-job_159757-607.jpg?size=626&ext=jpg",
		"https://img.freepik.com/premium-vector/deadline-time-late-work-busy-clock-pressure-countdown-lack-efficiency-watch_159757-684.jpg?size=626&ext=jpg",
		"https://img.freepik.com/premium-vector/turn-back-time-with-magnet-alarm-clock-return-timer-change-future-business_159757-685.jpg?w=360",
		"https://us.123rf.com/450wm/eamesbot/eamesbot2011/eamesbot201100051/159117873-lack-of-time-or-running-out-of-time-countdown-for-work-project-deadline-or-time-is-valuable-thing-in.jpg",
	}
)

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

type Schedule struct {
	UserID      string
	Appointment string
	Day         string
	Hours       string
}

func ScheduleAppointment(sc Schedule) error {
	var mtx sync.Mutex

	if err := verifyInputs(sc); err != nil {
		return err
	}

	mtx.Lock()
	defer mtx.Unlock()

	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("internal error, please try again")
	}
	defer f.Close()

	if _, err = f.WriteString(sc.UserID + "," + sc.Appointment + "," + sc.Day + " " + sc.Hours + "\n"); err != nil {
		return err
	}

	return nil
}

func verifyInputs(sc Schedule) error {
	if sc.UserID == "" || sc.Appointment == "" || sc.Day == "" || sc.Hours == "" {
		return errors.New("invalid values")
	}

	date := fmt.Sprintf("%s %s", sc.Day, sc.Hours)

	appointment, err := time.ParseInLocation("02/01/06 15:04", date, time.Local)
	if err != nil {
		return errors.New("invalid date format")
	}

	if time.Now().After(appointment) {
		return errors.New("selected period has passed")
	}

	return nil
}

func AlertAppointments(s *discordgo.Session) {
	for {
		f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDONLY, 0644)
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		buf := bytes.NewBuffer([]byte{})
		changed := false

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			record := strings.Split(line, ",")

			if len(record) >= 3 {
				date, err := time.ParseInLocation("02/01/06 15:04", record[2], time.Local)
				if err != nil {
					changed = true
					continue
				}

				now := time.Now()

				// date is in minute range, send appointment to user
				if date.After(now.Add(-1*time.Minute)) && date.Before(now) {
					go sendAlert(s, record[0], record[1], record[2])

					changed = true
					continue
				}

				// date has passed
				if date.Before(now) {
					changed = true
					continue
				}

				buf.WriteString(line + "\n")
			} else {
				changed = true
			}
		}

		if changed {
			err = os.WriteFile(fpath, buf.Bytes(), 0666)
			if err != nil {
				log.Println(err)
			}
		}

		f.Close()
		time.Sleep(10 * time.Second)
	}
}

func sendAlert(s *discordgo.Session, userID string, appointment string, date string) {
	c, err := s.UserChannelCreate(userID)
	if err != nil {
		return
	}

	s.ChannelMessageSendEmbed(c.ID, &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeGifv,
		Title:       "Reminder",
		Description: "There is your appointment!",
		Color:       colors[r.Intn(len(colors))],
		Image: &discordgo.MessageEmbedImage{
			URL: images[r.Intn(len(images))],
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Appointment",
				Value:  appointment,
				Inline: true,
			},
			{
				Name:   "Chosen date",
				Value:  date,
				Inline: true,
			},
		},
	})
}
