package main

import (
	"encoding/json"
	"internal/schedule"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func getIP(addr string) string {
	return addr[:strings.Index(addr, ":")]
}

func addScheduleHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var s schedule.Schedule
	json.NewDecoder(req.Body).Decode(&s)

	err := addSchedule(s)
	if err != nil {
		log.Printf("failed to insert schedule: %v", err)
		msg := Message{"failed to insert schedule"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	serverIP = getIP(req.RemoteAddr)
	msg := Message{"scheduled successfully"}
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

	err = removeSchedule(id)
	if err != nil {
		log.Printf("failed to delete schedule: %v", err)
		msg := Message{"failed to delete schedule"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	serverIP = getIP(req.RemoteAddr)
	msg := Message{"schedule removed successfully"}
	json.NewEncoder(w).Encode(msg)
}

func editScheduleHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var s schedule.Schedule
	json.NewDecoder(req.Body).Decode(&s)

	err := editSchedule(s)
	if err != nil {
		log.Printf("failed to update schedule: %v", err)
		msg := Message{"failed to update schedule"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	serverIP = getIP(req.RemoteAddr)
	msg := Message{"schedule edited successfully"}
	json.NewEncoder(w).Encode(msg)
}

func activateSettingHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	interval, err := strconv.Atoi(ps.ByName("interval"))
	if err != nil {
		log.Printf("invalid interval: %v", err)
		msg := Message{"invalid interval"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}

	settings[id] = interval

	serverIP = getIP(req.RemoteAddr)
	msg := Message{"setting activated successfully"}
	json.NewEncoder(w).Encode(msg)
}

func deactivateSettingHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if _, ok := settings[id]; ok {
		delete(settings, id)
	}

	serverIP = getIP(req.RemoteAddr)
	msg := Message{"setting activated successfully"}
	json.NewEncoder(w).Encode(msg)
}

func registerRouter(port string, r *httprouter.Router) {
	r.POST("/schedule/add", addScheduleHandler)
	r.POST("/schedule/remove/:id", removeScheduleHandler)
	r.POST("/schedule/edit", editScheduleHandler)

	r.POST("/setting/activate/:id/:interval", activateSettingHandler)
	r.POST("/setting/deactivate/:id", deactivateSettingHandler)

	log.Printf("Up and running at port=%s\n", port)
	http.ListenAndServe(port, r)
}
