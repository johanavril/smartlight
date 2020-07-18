package schedule

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Schedule struct {
	ID     int    `json:"id" db:"id"`
	LampID string `json:"lamp_id" db:"lamp_id"`
	Name   string `json:"name" db:"name"`
	Time   string `json:"time" db:"time"`
	TurnOn bool   `json:"turn_on" db:"turn_on"`
}

func GetSchedule(db *sql.DB, id int) (Schedule, error) {
	var s Schedule
	query := "SELECT id, lamp_id, name, time, turn_on FROM schedules WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&s.ID, &s.LampID, &s.Name, &s.Time, &s.TurnOn)
	if err != nil {
		log.Printf("failed to get schedule id=%d", id)
		return s, err
	}

	return s, nil
}

func GetLampSchedules(db *sql.DB, lampID string) ([]Schedule, error) {
	schedules := []Schedule{}
	query := "SELECT id, lamp_id, name, time, turn_on FROM schedules WHERE lamp_id = $1"
	rows, err := db.Query(query, lampID)
	if err != nil {
		log.Printf("failed to get schedule for lamp_id=%s", lampID)
		return nil, err
	}

	for rows.Next() {
		var s Schedule
		err = rows.Scan(&s.ID, &s.LampID, &s.Name, &s.Time, &s.TurnOn)
		if err != nil {
			break
		}
		schedules = append(schedules, s)
	}

	if err := rows.Err(); err != nil {
		log.Printf("failed to read schedule for lamp_id=%s", lampID)
		return nil, err
	}

	return schedules, nil
}

func GetAllSchedules(db *sql.DB) ([]Schedule, error) {
	schedules := []Schedule{}
	query := "SELECT id, lamp_id, name, time, turn_on FROM schedules"
	rows, err := db.Query(query)
	if err != nil {
		log.Print("failed to get all schedules")
		return nil, err
	}

	for rows.Next() {
		var s Schedule
		err = rows.Scan(&s.ID, &s.LampID, &s.Name, &s.Time, &s.TurnOn)
		if err != nil {
			break
		}
		schedules = append(schedules, s)
	}

	if err := rows.Err(); err != nil {
		log.Print("failed to read all schedules")
		return nil, err
	}

	return schedules, nil
}

func InsertSchedule(db *sql.DB, schedule Schedule) error {
	query := "INSERT INTO schedules(lamp_id, name, time, turn_on) VALUES($1, $2, $3, $4)"
	_, err := db.Exec(query, schedule.LampID, schedule.Name, schedule.Time, schedule.TurnOn)

	if err != nil {
		log.Printf("failed to insert schedule %+v", schedule)
		return err
	}

	return nil
}

func DeleteSchedule(db *sql.DB, id int) error {
	query := "DELETE FROM schedules WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		log.Printf("failed to delete schedule id=%d", id)
		return err
	}

	return nil
}

func UpdateSchedule(db *sql.DB, schedule Schedule) error {
	query := "UPDATE schedules SET name = $1, time = $2, turn_on = $3 WHERE id = $4"
	_, err := db.Exec(query, schedule.Name, schedule.Time, schedule.TurnOn, schedule.ID)
	if err != nil {
		log.Printf("failed to update schedule id=%d %+v", schedule.ID, schedule)
		return err
	}

	return nil
}

func ScheduleBot(db *sql.DB, s Schedule, botIDAddress map[string]string) error {
	if s.ID == 0 {
		schedules, err := GetLampSchedules(db, s.LampID)
		if err != nil {
			log.Printf("schedules not found: %v", err)
			return err
		}
		s = schedules[len(schedules)-1]
	}

	body, err := json.Marshal(s)
	if err != nil {
		log.Printf("failed to marshal schedule: %v", err)
		return err
	}

	botID := s.LampID[:strings.Index(s.LampID, ":")]
	ip, ok := botIDAddress[botID]
	if !ok {
		log.Printf("address not found for botID=%s", botID)
		return fmt.Errorf("bot address not found")
	}

	resp, err := http.Post(fmt.Sprintf("http://%s:10001/schedule/add", ip), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("failed to send schedule: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Print("failed to schedule bot")
		return fmt.Errorf("request to bot failed")
	}

	return nil
}

func UnscheduleBot(db *sql.DB, id int, botIDAddress map[string]string) error {
	s, err := GetSchedule(db, id)
	if err != nil {
		log.Printf("schedule not found: %v", err)
		return err
	}

	botID := s.LampID[:strings.Index(s.LampID, ":")]
	ip, ok := botIDAddress[botID]
	if !ok {
		log.Printf("address not found for botID=%s", botID)
		return fmt.Errorf("bot address not found")
	}

	resp, err := http.Post(fmt.Sprintf("http://%s:10001/schedule/remove/%d", ip, id), "application/json", nil)
	if err != nil {
		log.Printf("failed to remove schedule: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Print("failed to unschedule bot")
		return fmt.Errorf("request to bot failed")
	}

	return nil
}

func RescheduleBot(s Schedule, botIDAddress map[string]string) error {
	body, err := json.Marshal(s)
	if err != nil {
		log.Printf("failed to marshal schedule: %v", err)
		return err
	}

	botID := s.LampID[:strings.Index(s.LampID, ":")]
	ip, ok := botIDAddress[botID]
	if !ok {
		log.Printf("address not found for botID=%s", botID)
		return fmt.Errorf("bot address not found")
	}

	resp, err := http.Post(fmt.Sprintf("http://%s:10001/schedule/edit", ip), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("failed to update schedule: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Print("failed to reschedule bot")
		return fmt.Errorf("request to bot failed")
	}

	return nil
}
