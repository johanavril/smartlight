package main

import (
	"encoding/json"
	"internal/connection"
	"internal/usage"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getLampUsagesHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id := ps.ByName("lamp-id")

	usages, err := usage.GetLampUsages(connection.DB, id)
	if err != nil {
		log.Printf("failed to get usage for lamp id=%s from database: %v", id, err)
		msg := Message{"failed to get usage(s) from database"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	json.NewEncoder(w).Encode(usages)
}

func getAllUsagesHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	usages, err := usage.GetAllUsages(connection.DB)
	if err != nil {
		log.Printf("failed to get all usages from database: %v", err)
		msg := Message{"failed to get all usages from database"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	json.NewEncoder(w).Encode(usages)
}

func addUsageHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var u usage.Usage
	json.NewDecoder(req.Body).Decode(&u)

	err := usage.InsertUsage(connection.DB, u)
	if err != nil {
		log.Printf("failed to insert usage: %v", err)
		msg := Message{"failed to insert usage"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Message{"insert usage success"}
	json.NewEncoder(w).Encode(msg)
}
