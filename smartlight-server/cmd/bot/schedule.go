package main

import (
	"fmt"
	"internal/schedule"
	"internal/usage"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"gobot.io/x/gobot/drivers/gpio"
)

var scheduler *gocron.Scheduler
var schedules = make(map[int]*gocron.Job)

func task(r *gpio.RelayDriver, s schedule.Schedule) {
	if s.TurnOn {
		r.Off()
		executeSetting(r, s.LampID)
	} else {
		r.On()
	}

	var u usage.Usage
	u.LampID = s.LampID
	u.Time = time.Now()
	u.TurnOn = s.TurnOn
	storeUsage(u)
}

func addSchedule(s schedule.Schedule) error {
	if _, ok := schedules[s.ID]; ok {
		log.Printf("schedule already exists %+v", s)
		return nil
	}

	r, ok := relays[s.LampID]
	if !ok {
		return fmt.Errorf("lamp id=%s is not controlled by this bot", s.LampID)
	}

	job, err := scheduler.Every(1).Day().At(s.Time).Do(task, r, s)
	if err != nil {
		return fmt.Errorf("failed to register schedule %v: %v", s, err)
	}

	schedules[s.ID] = job

	return nil
}

func removeSchedule(id int) error {
	s, ok := schedules[id]
	if !ok {
		return fmt.Errorf("schedule id=%d has not been scheduled", id)
	}
	scheduler.RemoveByReference(s)

	return nil
}

func editSchedule(s schedule.Schedule) error {
	r, ok := relays[s.LampID]
	if !ok {
		return fmt.Errorf("lamp id=%s is not controlled by this bot", s.LampID)
	}

	if err := removeSchedule(s.ID); err != nil {
		return err
	}

	job, err := scheduler.Every(1).Day().At(s.Time).Do(task, r, s)
	if err != nil {
		return fmt.Errorf("failed to edit schedule %v: %v", s, err)
	}

	schedules[s.ID] = job

	return nil
}
