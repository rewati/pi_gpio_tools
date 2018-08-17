package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const blinkHelp = "Blink: Will start continuously toggle state on pin. Need pin number and blink lag in seconds. Example:\nsudo pinion blink 7 1"
const toggleHelp = "Toggle: Will toggle stage from on to off or off to on. Need pin number. Example:\nsudo toggle 7"

const help = "Usage:" + "\n1. " + blinkHelp + "\n2. " + toggleHelp

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Invalid command. \n", help)
		os.Exit(-1)
	}

	command := os.Args[1]
	switch command {
	case "blink":
		createBlinkRobot()
	case "toggle":
		createToggleRobot()
	default:
		fmt.Println("Invalid command. \n", help)
		os.Exit(-1)
	}
}

func createBlinkRobot() {
	if len(os.Args) < 4 {
		fmt.Println("Invalid command. \n", blinkHelp)
		os.Exit(-1)
	}
	pin := getPin()
	s := os.Args[3]

	sec, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("%v has to be int.\n", s)
		fmt.Println("Invalid command. ", help)
		os.Exit(-1)
	}

	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, pin)

	work := func() {
		gobot.Every(time.Duration(sec)*time.Second, func() {
			led.Toggle()
		})
	}

	gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	).Start()
}

func createToggleRobot() {
	if len(os.Args) < 3 {
		fmt.Println("Invalid command. \n", toggleHelp)
		os.Exit(-1)
	}
	pin := getPin()
	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, pin)
	work := func() {
		fmt.Printf("Starting to toggle state of pin: %v\n", pin)
		led.Toggle()
		fmt.Printf("Toggled state of pin: %v\n", pin)
	}
	robot := gobot.NewRobot("toggleBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)
	robot.Work()
}

func getPin() string {
	pin := os.Args[2]
	if _, err := strconv.Atoi(pin); err != nil {
		fmt.Printf("%v has to be int.\n", pin)
		fmt.Println("Invalid command. ", help)
		os.Exit(-1)
	}
	return pin
}
