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
	history   map[string]int
}

func chooseSound(dirname string, history map[string]int) string {
	rand.Seed(time.Now().Unix())

	fileList := []string{}
	_ = filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "mp3") {
			fileList = append(fileList, path)
		}
		return nil
	})

	filtered := []string{}
	for _,file := range fileList {
		seen := history[file]
		if seen==0 {
			filtered = append(filtered, file)
		}
	}

	file := filtered[rand.Intn(len(filtered))]
	if len(filtered) == 1 {
		for k:= range history {
			history[k]=0
		}
	}
	history[file] = 1

	return file
}

func makeButtonPushHandler(buttonName string) func(data interface{}) {
	name := buttonName
	return func(data interface{}) {
		button := buttons[name]
		fmt.Println(button.name + " pushed")
		// fmt.Println(button.history)

		sound := chooseSound(button.soundPath, button.history)

		command := exec.Command("play", sound)

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
		button.history[sound] = 1
	}
}

func makeButtonReleaseHandler(buttonName string) func(data interface{}) {
	name := buttonName
	return func(data interface{}) {
		button := buttons[name]
		p := <-button.process
		if p != nil {
			fmt.Println(p)
			p.Kill()
			button.led.Off()
		}
	}
}

var buttons map[string]ArcadeButton

func main() {

	r := raspi.NewRaspiAdaptor("raspi")

	buttons = map[string]ArcadeButton{
		"red":    {"red", gpio.NewLedDriver(r, "led-red", "3"), gpio.NewButtonDriver(r, "button-red", "15"), "./sounds/tragic", make(chan *os.Process), make(map[string]int)},
		"green":  {"green", gpio.NewLedDriver(r, "led-green", "5"), gpio.NewButtonDriver(r, "button-green", "19"), "./sounds/tagueule", make(chan *os.Process), make(map[string]int)},
		"yellow": {"yellow", gpio.NewLedDriver(r, "led-yellow", "7"), gpio.NewButtonDriver(r, "button-yellow", "21"), "./sounds/wtf", make(chan *os.Process), make(map[string]int)},
		"blue":   {"blue", gpio.NewLedDriver(r, "led-blue", "11"), gpio.NewButtonDriver(r, "button-blue", "23"), "./sounds/yeah", make(chan *os.Process), make(map[string]int)},
		"white":  {"white", gpio.NewLedDriver(r, "led-white", "13"), gpio.NewButtonDriver(r, "button-white", "12"), "./sounds/slap", make(chan *os.Process), make(map[string]int)},
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
			gobot.On(b.button.Event("push"), makeButtonPushHandler(b.name))
			gobot.On(b.button.Event("release"), makeButtonReleaseHandler(b.name))
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
