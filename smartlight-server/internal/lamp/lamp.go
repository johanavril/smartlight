package lamp

import (
	"database/sql"
	"log"
)

type Lamp struct {
	ID         int    `json:"id,omitempty" db:"id"`
	Name       string `json:"name" db:"name"`
	TotalLamp  int    `json:"total_lamp" db:"total_lamp"`
	TotalPower int    `json:"total_power" db:"total_power"`
	ImageName  string `json:"image_name" db:"image_name"`
}

func GetLamp(db *sql.DB, id int) (Lamp, error) {
	var l Lamp
	query := "SELECT id, name, total_lamp, total_power, image_name FROM lamps WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&l.ID, &l.Name, &l.TotalLamp, &l.TotalPower, &l.ImageName)
	if err != nil {
		log.Printf("failed to get lamp id=%d", id)
		return l, err
	}

	return l, nil
}

func GetLamps(db *sql.DB) ([]Lamp, error) {
	var lamps []Lamp
	query := "SELECT id, name, total_lamp, total_power, image_name FROM lamps"
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Print("failed to query lamps")
		return lamps, err
	}

	for rows.Next() {
		var l Lamp
		err = rows.Scan(&l.ID, &l.Name, &l.TotalLamp, &l.TotalPower, &l.ImageName)
		if err != nil {
			break
		}
		lamps = append(lamps, l)
	}

	if err := rows.Err(); err != nil {
		log.Print("failed to read lamps")
		return lamps, err
	}

	return lamps, nil
}

func InsertLamp(db *sql.DB, lamp Lamp) error {
	query := "INSERT INTO lamps(name, total_lamp, total_power, image_name) VALUES($1, $2, $3, $4)"
	_, err := db.Exec(query, lamp.Name, lamp.TotalLamp, lamp.TotalPower, lamp.ImageName)
	if err != nil {
		log.Printf("failed to insert lamp %+v", lamp)
		return err
	}

	return nil
}
