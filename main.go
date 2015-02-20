package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

type ArcadeButton struct {
	name      string
	led       *gpio.LedDriver
	button    *gpio.ButtonDriver
	soundPath string
	process   chan *os.Process
}

func chooseSound(dirname string) string {
	rand.Seed(time.Now().Unix())

	fileList := []string{}
	_ = filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "mp3") {
			fileList = append(fileList, path)
		}
		return nil
	})

	file := fileList[rand.Intn(len(fileList))]

	return file
}

func makeButtonPushHandler(b ArcadeButton) func(data interface{}) {
	button := b
	return func(data interface{}) {
		fmt.Println(button.name + " pushed")

		sound := chooseSound(button.soundPath)

		command := exec.Command("play", chooseSound(button.soundPath))

		button.led.On()
		if strings.Contains(sound, "pushmode") {
			command.Start()
			fmt.Println(command.Process)
			button.process <- command.Process
		} else {
			command.Run()
			button.led.Off()
			button.process <- nil
		}
	}
}

func makeButtonReleaseHandler(b ArcadeButton) func(data interface{}) {
	button := b
	return func(data interface{}) {
		p := <-button.process
		if p != nil {
			fmt.Println(p)
			p.Kill()
			button.led.Off()
		}
	}
}

func main() {

	r := raspi.NewRaspiAdaptor("raspi")

	var buttons = map[string]ArcadeButton{
		"red":    {"red", gpio.NewLedDriver(r, "led-red", "3"), gpio.NewButtonDriver(r, "button-red", "15"), "./sounds/tragic", make(chan *os.Process)},
		"green":  {"green", gpio.NewLedDriver(r, "led-green", "5"), gpio.NewButtonDriver(r, "button-green", "19"), "./sounds/tagueule", make(chan *os.Process)},
		"yellow": {"yellow", gpio.NewLedDriver(r, "led-yellow", "7"), gpio.NewButtonDriver(r, "button-yellow", "21"), "./sounds/wtf", make(chan *os.Process)},
		"blue":   {"blue", gpio.NewLedDriver(r, "led-blue", "11"), gpio.NewButtonDriver(r, "button-blue", "23"), "./sounds/yeah", make(chan *os.Process)},
		"white":  {"white", gpio.NewLedDriver(r, "led-white", "13"), gpio.NewButtonDriver(r, "button-white", "12"), "./sounds/slap", make(chan *os.Process)},
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
		for _, b := range buttons {
			gobot.On(b.button.Event("push"), makeButtonPushHandler(b))
			gobot.On(b.button.Event("release"), makeButtonReleaseHandler(b))
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
