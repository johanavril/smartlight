package main

import (
	"encoding/json"
	"internal/connection"
	"internal/lamp"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getLampHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id := ps.ByName("id")

	l, err := lamp.GetLamp(connection.DB, id)
	if err != nil {
		log.Printf("failed to get lamp id=%s from database: %v", id, err)
		msg := Message{"failed to get data from database"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	json.NewEncoder(w).Encode(l)
}

func getLampsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	lamps, err := lamp.GetLamps(connection.DB)
	if err != nil {
		log.Printf("failed to get all lamp from database: %v", err)
		msg := Message{"failed to get data from database"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	json.NewEncoder(w).Encode(lamps)
}

func addLampHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var l lamp.Lamp
	json.NewDecoder(req.Body).Decode(&l)

	err := lamp.InsertLamp(connection.DB, l)
	if err != nil {
		log.Printf("failed to insert lamp: %v", err)
		msg := Message{"failed to insert lamp"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Message{"insert lamp success"}
	json.NewEncoder(w).Encode(msg)
}
