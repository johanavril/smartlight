package usage

import (
	"database/sql"
	"log"
	"time"
)

type Usage struct {
	ID     int       `json:"id" db:"id"`
	LampID int       `json:"lamp_id" db:"lamp_id"`
	Time   time.Time `json:"time" db:"time"`
	TurnOn bool      `json:"turn_on" db:"turn_on"`
}

func GetLampUsages(db *sql.DB, lampID int) ([]Usage, error) {
	usages := []Usage{}
	query := "SELECT id, lamp_id, time, turn_on FROM usages WHERE lamp_id = $1"
	rows, err := db.Query(query, lampID)
	if err != nil {
		log.Printf("failed to get usage for lamp_id=%d", lampID)
		return nil, err
	}

	for rows.Next() {
		var u Usage
		err = rows.Scan(&u.ID, &u.LampID, &u.Time, &u.TurnOn)
		if err != nil {
			break
		}
		usages = append(usages, u)
	}

	if err := rows.Err(); err != nil {
		log.Printf("failed to read usage for lamp_id=%d", lampID)
		return nil, err
	}

	return usages, nil
}

func GetAllUsages(db *sql.DB) ([]Usage, error) {
	usages := []Usage{}
	query := "SELECT id, lamp_id, time, turn_on FROM usages"
	rows, err := db.Query(query)
	if err != nil {
		log.Print("failed to get all usages")
		return nil, err
	}

	for rows.Next() {
		var u Usage
		err = rows.Scan(&u.ID, &u.LampID, &u.Time, &u.TurnOn)
		if err != nil {
			break
		}
		usages = append(usages, u)
	}

	if err := rows.Err(); err != nil {
		log.Print("failed to read all usages")
		return nil, err
	}

	return usages, nil
}

func InsertUsage(db *sql.DB, usage Usage) error {
	query := "INSERT INTO usages(lamp_id, time, turn_on) VALUES($1, $2, $3)"
	_, err := db.Exec(query, usage.LampID, usage.Time, usage.TurnOn)

	if err != nil {
		log.Printf("failed to insert usage %+v", usage)
		return err
	}

	return nil
}
