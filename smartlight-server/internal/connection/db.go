package connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() error {
	username := os.Getenv("SMARTLIGHT_USERNAME")
	password := os.Getenv("SMARTLIGHT_PASSWORD")
	dbname := os.Getenv("SMARTLIGHT_DB")
	sslmode := os.Getenv("SMARTLIGHT_SSLMODE")
	conn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", username, password, dbname, sslmode)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Print("failed to create database connection")
		return err
	}

	if err = db.Ping(); err != nil {
		log.Print("failed to communicate with the database")
		return err
	}
	DB = db

	return nil
}
