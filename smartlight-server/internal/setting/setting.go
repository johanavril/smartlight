package setting

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type Setting struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Interval int    `json:"interval" db:"interval"`
	Checked  bool   `json:"checked" db:"checked"`
}

func GetSetting(db *sql.DB, id string) (Setting, error) {
	var s Setting
	query := "SELECT id, name, interval, checked FROM settings WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&s.ID, &s.Name, &s.Interval, &s.Checked)
	if err != nil {
		log.Printf("failed to get setting id=%d", id)
		return s, err
	}

	return s, nil
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

func SetSetting(s Setting, botIDAddress map[string]string) {
	for _, ip := range botIDAddress {
		var URL string
		if s.Checked {
			URL = fmt.Sprintf("http://%s:10001/setting/activate/%s/%d", ip, s.ID, s.Interval)
		} else {
			URL = fmt.Sprintf("http://%s:10001/setting/deactivate/%s", ip, s.ID)
		}

		resp, err := http.Post(URL, "application/json", nil)
		if err != nil {
			log.Printf("failed to set setting=%v: %v", s, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Print("failed to set setting bot")
		}
	}
}
