package main

import (
	"encoding/json"
	"internal/connection"
	"internal/schedule"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func getLampSchedulesHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	id := ps.ByName("lamp-id")

	schedules, err := schedule.GetLampSchedules(connection.DB, id)
	if err != nil {
		log.Printf("failed to get schedule for lamp id=%s from database: %v", id, err)
		msg := Message{"failed to get schedule(s) from database"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	json.NewEncoder(w).Encode(schedules)
}

func getAllSchedulesHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	schedules, err := schedule.GetAllSchedules(connection.DB)
	if err != nil {
		log.Printf("failed to get all schedule from database: %v", err)
		msg := Message{"failed to get all schedules from database"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	json.NewEncoder(w).Encode(schedules)
}

func addScheduleHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var s schedule.Schedule
	json.NewDecoder(req.Body).Decode(&s)

	err := schedule.InsertSchedule(connection.DB, s)
	if err != nil {
		log.Printf("failed to insert schedule: %v", err)
		msg := Message{"failed to insert schedule"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Message{"insert schedule success"}
	json.NewEncoder(w).Encode(msg)
}

func removeScheduleHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Printf("invalid schedule id: %v", err)
		msg := Message{"invalid schedule id"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}

	err = schedule.DeleteSchedule(connection.DB, id)
	if err != nil {
		log.Printf("failed to delete schedule: %v", err)
		msg := Message{"failed to delete schedule"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Message{"delete schedule success"}
	json.NewEncoder(w).Encode(msg)
}

func editScheduleHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var s schedule.Schedule
	json.NewDecoder(req.Body).Decode(&s)

	err := schedule.UpdateSchedule(connection.DB, s)
	if err != nil {
		log.Printf("failed to update schedule: %v", err)
		msg := Message{"failed to update schedule"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	msg := Message{"update schedule success"}
	json.NewEncoder(w).Encode(msg)
}
