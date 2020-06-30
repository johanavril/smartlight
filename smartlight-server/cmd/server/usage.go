package main

import (
	"encoding/json"
	"internal/connection"
	"internal/usage"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func getLampUsagesHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(ps.ByName("lamp-id"))
	if err != nil {
		log.Printf("invalid lamp id: %v", err)
		msg := Message{"invalid lamp id"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	usages, err := usage.GetLampUsages(connection.DB, id)
	if err != nil {
		log.Printf("failed to get usage for lamp id=%d from database: %v", id, err)
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
