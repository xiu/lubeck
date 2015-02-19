package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

type ArcadeButton struct {
	led       *gpio.LedDriver
	button    *gpio.ButtonDriver
	soundPath string
}

func chooseSound(dirname string) string {
	rand.Seed(time.Now().Unix())

	files, _ := ioutil.ReadDir(dirname)

	file := files[rand.Intn(len(files))]

	return dirname + "/" + file.Name()
}

func main() {

	r := raspi.NewRaspiAdaptor("raspi")

	var buttons = map[string]ArcadeButton{
		"red":    {gpio.NewLedDriver(r, "led-red", "3"), gpio.NewButtonDriver(r, "button-red", "15"), "./sounds/tragic"},
		"green":  {gpio.NewLedDriver(r, "led-green", "5"), gpio.NewButtonDriver(r, "button-green", "19"), "./sounds/tagueule"},
		"yellow": {gpio.NewLedDriver(r, "led-yellow", "7"), gpio.NewButtonDriver(r, "button-yellow", "21"), "./sounds/wtf"},
		"blue":   {gpio.NewLedDriver(r, "led-blue", "11"), gpio.NewButtonDriver(r, "button-blue", "23"), "./sounds/yeah"},
		"white":  {gpio.NewLedDriver(r, "led-white", "13"), gpio.NewButtonDriver(r, "button-white", "12"), "./sounds/slap"},
	}

	gbot := gobot.NewGobot()

	allOff := func() {
		for _, b := range buttons {
			b.led.Off()
		}
	}

	allOn := func() {
		for _, b := range buttons {
			b.led.On()
			time.Sleep(100 * time.Millisecond)
		}
	}

	allOff()
	allOn()
	allOff()

	work := func() {
		for name, b := range buttons {
			gobot.On(b.button.Event("push"), func(data interface{}) {
				fmt.Println(name + " pushed")
				command := exec.Command("play", chooseSound(b.soundPath))
				b.led.On()
				command.Run()
				b.led.Off()
			})
		}
	}

	robot := gobot.NewRobot("Lubeck",
		[]gobot.Connection{r},
		[]gobot.Device{
			buttons["red"].button, buttons["green"].button, buttons["yellow"].button, buttons["blue"].button, buttons["white"].button,
			buttons["red"].led, buttons["green"].led, buttons["yellow"].led, buttons["blue"].led, buttons["white"].led,
		},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
