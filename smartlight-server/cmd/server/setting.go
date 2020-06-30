package main

import (
	"encoding/json"
	"internal/connection"
	"internal/setting"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getAllSettingsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	settings, err := setting.GetAllSettings(connection.DB)
	if err != nil {
		log.Printf("failed to get all setting from database: %v", err)
		msg := Message{"failed to get all settings from database"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	json.NewEncoder(w).Encode(settings)
}

func addSettingHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var s setting.Setting
	json.NewDecoder(req.Body).Decode(&s)

	err := setting.InsertSetting(connection.DB, s)
	if err != nil {
		log.Printf("failed to insert setting: %v", err)
		msg := Message{"failed to insert setting"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Message{"insert setting success"}
	json.NewEncoder(w).Encode(msg)
}

func removeSettingHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	err := setting.DeleteSetting(connection.DB, ps.ByName("id"))
	if err != nil {
		log.Printf("failed to delete setting: %v", err)
		msg := Message{"failed to delete setting"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Message{"delete setting success"}
	json.NewEncoder(w).Encode(msg)
}

func editSettingHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var s setting.Setting
	json.NewDecoder(req.Body).Decode(&s)

	err := setting.UpdateSetting(connection.DB, s)
	if err != nil {
		log.Printf("failed to update setting: %v", err)
		msg := Message{"failed to update setting"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Message{"update setting success"}
	json.NewEncoder(w).Encode(msg)
}