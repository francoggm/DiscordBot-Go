package schedule

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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
	if err := verifyInputs(sc); err != nil {
		return err
	}

	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("internal error, please try again")
	}
	defer f.Close()

	if err = verifyAppointmentExists(sc); err != nil {
		return err
	}

	if err = writeAppointment(f, sc); err != nil {
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

func verifyAppointmentExists(sc Schedule) error {
	record := fmt.Sprintf("%s,%s,%s %s", sc.UserID, sc.Appointment, sc.Day, sc.Hours)

	content, err := os.ReadFile(fpath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, record) {
			return errors.New("appointment already exists")
		}
	}

	return nil
}

func writeAppointment(f *os.File, sc Schedule) error {
	record := fmt.Sprintf("%s,%s,%s %s", sc.UserID, sc.Appointment, sc.Day, sc.Hours)

	_, err := f.WriteString(record + "\n")
	return err
}

func AlertAppointments(s *discordgo.Session) {
	for {
		f, err := os.Open(fpath)
		if err != nil {
			os.WriteFile(fpath, []byte{}, 0666)
		}

		var bs []byte
		buf := bytes.NewBuffer(bs)

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			record := strings.Split(line, ",")

			if len(record) >= 3 {
				date, err := time.ParseInLocation("02/01/06 15:04", record[2], time.Local)
				if err != nil {
					continue
				}

				if date.After(time.Now().Add(-1*time.Minute)) && date.Before(time.Now()) {
					sendAlert(record[0], record[1])
					continue
				}

				if date.Before(time.Now()) {
					continue
				}

				buf.WriteString(line)
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

func sendAlert(userID string, appointment string) {
	// send to user
	fmt.Println("Appointment=" + appointment + " to " + userID)
}
