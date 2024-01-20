package schedule

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	fpath = "schedule/schedule.txt"
)

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

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			record := strings.Split(line, ",")

			if len(record) >= 3 {
				date, err := time.ParseInLocation("02/01/06 15:04", record[2], time.Local)
				if err != nil {
					continue
				}

				now := time.Now()

				// date is in minute range, send appointment to user
				if date.After(now.Add(-1*time.Minute)) && date.Before(now) {
					go func(session *discordgo.Session, userID string, appointment string) {
						c, err := session.UserChannelCreate(userID)
						if err != nil {
							return
						}
						session.ChannelMessageSend(c.ID, appointment)
					}(s, record[0], record[1])

					continue
				}

				// date has passed
				if date.Before(now) {
					continue
				}

				buf.WriteString(line + "\n")
			}
		}

		err = os.WriteFile(fpath, buf.Bytes(), 0666)
		if err != nil {
			log.Fatal(err)
		}

		f.Close()
		time.Sleep(10 * time.Second)
	}
}
