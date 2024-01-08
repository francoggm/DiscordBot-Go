package schedule

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	scheduleFile = "schedule/schedule.txt"
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

	f, err := os.OpenFile(scheduleFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("internal error, please try again")
	}
	defer f.Close()

	if err = verifyAppointmentExists(f, sc); err != nil {
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

	appointmentDay, err := time.Parse("02/01/06 15:04", date)
	if err != nil {
		return errors.New("invalid date format")
	}

	// TODO: Fix comparation dates
	if time.Now().Before(appointmentDay) {
		return errors.New("selected period has passed")
	}

	return nil
}

func verifyAppointmentExists(f *os.File, sc Schedule) error {
	record := fmt.Sprintf("%s,%s,%s %s", sc.UserID, sc.Appointment, sc.Day, sc.Hours)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), record) {
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
