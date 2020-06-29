package setting

import (
	"database/sql"
	"log"
)

type Setting struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Interval int    `json:"interval" db:"interval"`
	Checked  bool   `json:"checked" db:"checked"`
}

func GetAllSettings(db *sql.DB) ([]Setting, error) {
	settings := []Setting{}
	query := "SELECT id, name, interval, checked FROM settings"
	rows, err := db.Query(query)
	if err != nil {
		log.Print("failed to get all settings")
		return nil, err
	}

	for rows.Next() {
		var s Setting
		err = rows.Scan(&s.ID, &s.Name, &s.Interval, &s.Checked)
		if err != nil {
			break
		}
		settings = append(settings, s)
	}

	if err := rows.Err(); err != nil {
		log.Print("failed to read all settings")
		return nil, err
	}

	return settings, nil
}

func InsertSetting(db *sql.DB, setting Setting) error {
	query := "INSERT INTO settings(id, name, interval, checked) VALUES($1, $2, $3, $4)"
	_, err := db.Exec(query, setting.ID, setting.Name, setting.Interval, setting.Checked)
	if err != nil {
		log.Printf("failed to insert setting %+v", setting)
		return err
	}

	return nil
}

func DeleteSetting(db *sql.DB, id string) error {
	query := "DELETE FROM settings WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		log.Printf("failed to delete setting id=%s", id)
		return err
	}

	return nil
}

func UpdateSetting(db *sql.DB, setting Setting) error {
	query := "UPDATE settings SET name = $1, interval = $2, checked = $3 WHERE id = $4"
	_, err := db.Exec(query, setting.Name, setting.Interval, setting.Checked, setting.ID)
	if err != nil {
		log.Printf("failed to update setting id=%s %+v", setting.ID, setting)
		return err
	}

	return nil
}
