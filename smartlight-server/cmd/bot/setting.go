package main

import (
	"internal/usage"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

var settings = make(map[string]int)

func executeSetting(r *gpio.RelayDriver, lampID string) {
	for _, v := range settings {
		gobot.After(time.Duration(v)*time.Minute, func() {
			r.On()

			var u usage.Usage
			u.LampID = lampID
			u.Time = time.Now()
			u.TurnOn = false
			storeUsage(u)
		})
	}
}
