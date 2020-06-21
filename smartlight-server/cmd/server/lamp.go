package main

import (
	"encoding/json"
	"internal/connection"
	"internal/lamp"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func getLampHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("invalid lamp id: %v", err)
		msg := Message{"invalid lamp id"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	l, err := lamp.GetLamp(connection.DB, id)
	if err != nil {
		log.Printf("failed to get lamp id=%d from database: %v", id, err)
		msg := Message{"failed to get data from database"}
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
		json.NewEncoder(w).Encode(msg)
		return
	}

	if err != nil {
		log.Printf("failed to marshal lamps: %v", err)
		msg := Message{"failed to build response"}
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
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Message{"insert lamp success"}
	json.NewEncoder(w).Encode(msg)
}
