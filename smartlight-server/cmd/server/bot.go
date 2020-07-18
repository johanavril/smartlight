package main

import (
	"internal/connection"
	"internal/schedule"
	"internal/setting"
	"log"
	"time"

	"github.com/go-co-op/gocron"
)

var botIDAddress = map[string]string{
	"1": "192.168.0.50",
}

func syncSchedules() {
	schedules, err := schedule.GetAllSchedules(connection.DB)
	if err != nil {
		log.Panicf("failed to get schedules: %v", err)
	}

	for _, s := range schedules {
		schedule.ScheduleBot(connection.DB, s, botIDAddress)
	}
}

func syncSettings() {
	settings, err := setting.GetAllSettings(connection.DB)
	if err != nil {
		log.Panicf("failed to get settings: %v", err)
	}

	for _, s := range settings {
		if !s.Checked {
			continue
		}
		setting.SetSetting(s, botIDAddress)
	}
}

func autoSync() {
	syncSchedules()
	syncSettings()

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(5).Minutes().Do(func() {
		syncSchedules()
		syncSettings()
	})

	scheduler.StartAsync()
}
