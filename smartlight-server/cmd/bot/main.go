package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/julienschmidt/httprouter"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

var serverIP string
var relays = make(map[string]*gpio.RelayDriver)

type Message struct {
	Message string `json:"message"`
}

func main() {
	serverIP = "192.168.0.100" // default, most likely useless because serverIP use DHCP
	router := httprouter.New()
	scheduler = gocron.NewScheduler(time.UTC)

	botID := os.Getenv("BOT_ID")
	botName := fmt.Sprintf("Lamp Bot %s", botID)
	port := fmt.Sprintf(":%d", 10001)

	r := raspi.NewAdaptor()
	in1 := gpio.NewRelayDriver(r, os.Getenv("RELAY_IN1"))
	in2 := gpio.NewRelayDriver(r, os.Getenv("RELAY_IN2"))
	in3 := gpio.NewRelayDriver(r, os.Getenv("RELAY_IN3"))
	in4 := gpio.NewRelayDriver(r, os.Getenv("RELAY_IN4"))

	// format lamp_id = BOT_ID:RELAY_CHANNEL
	// change the RELAY_CHANNEL value according to the intended channel to be controlled
	relays[fmt.Sprintf("%s:%d", botID, 1)] = in1
	relays[fmt.Sprintf("%s:%d", botID, 2)] = in2
	relays[fmt.Sprintf("%s:%d", botID, 3)] = in3
	relays[fmt.Sprintf("%s:%d", botID, 4)] = in4

	work := func() {
		gobot.Every(5*time.Minute, func() {
			sendUsages()
		})
		scheduler.StartAsync()
		registerRouter(port, router)
	}

	robot := gobot.NewRobot(botName,
		[]gobot.Connection{r},
		[]gobot.Device{in1, in2, in3, in4},
		work,
	)

	robot.Start()
	robot.Devices().Start()
}
