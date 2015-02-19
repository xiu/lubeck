package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

func main() {
	gbot := gobot.NewGobot()

	r := raspi.NewRaspiAdaptor("raspi")
	ledred := gpio.NewLedDriver(r, "led-red", "3")
	ledgreen := gpio.NewLedDriver(r, "led-green", "5")
	ledyellow := gpio.NewLedDriver(r, "led-yellow", "7")
	ledblue := gpio.NewLedDriver(r, "led-blue", "11")
	ledwhite := gpio.NewLedDriver(r, "led-white", "13")
	buttonred := gpio.NewButtonDriver(r, "button-red", "15")
	buttongreen := gpio.NewButtonDriver(r, "button-green", "19")
        buttonyellow := gpio.NewButtonDriver(r, "button-yellow", "21")
        buttonblue := gpio.NewButtonDriver(r, "button-blue", "23")
        buttonwhite := gpio.NewButtonDriver(r, "button-white", "12")

	ledred.Off()
	ledgreen.Off()
	ledyellow.Off()
	ledblue.Off()
	ledwhite.Off()

	ledred.On()
	time.Sleep(100 * time.Millisecond)
	ledgreen.On()
	time.Sleep(100 * time.Millisecond)
	ledyellow.On()
	time.Sleep(100 * time.Millisecond)
	ledblue.On()
	time.Sleep(100 * time.Millisecond)
	ledwhite.On()
	time.Sleep(100 * time.Millisecond)

        ledred.Off()
        ledgreen.Off()
        ledyellow.Off()
        ledblue.Off()
        ledwhite.Off()
	

	work := func() {
		gobot.On(buttonred.Event("push"), func(data interface{}) {
			inception := exec.Command("play", "/root/inception.mp3")
			fmt.Println("PUSHED RED")
			ledred.On()
			inception.Run()
			ledred.Off()
		})

                gobot.On(buttongreen.Event("push"), func(data interface{}) {
		        blabla := exec.Command("play", "/root/blabla.mp3")
                        fmt.Println("PUSHED GREEN")
                        ledgreen.On()
                        blabla.Run()
                        ledgreen.Off()
                })

		gobot.On(buttonyellow.Event("push"), func(data interface{}) {
	                cotcot := exec.Command("play", "/root/cotcot.mp3")
			fmt.Println("PUSHED YELLOW")
			ledyellow.On()
			cotcot.Run()
			ledyellow.Off()
		})	
		
                gobot.On(buttonblue.Event("push"), func(data interface{}) {
	                toilets := exec.Command("play", "/root/toilets.mp3")
                        fmt.Println("PUSHED BLUE")
                        ledblue.On()
                        toilets.Run()
                        ledblue.Off()
                })

                gobot.On(buttonwhite.Event("push"), func(data interface{}) {
	                whip := exec.Command("play", "/root/whip.mp3")
                        fmt.Println("PUSHED WHITE")
                        ledwhite.On()
                        whip.Run()
                        ledwhite.Off()
                })
	}

	robot := gobot.NewRobot("Lubeck",
		[]gobot.Connection{r},
		[]gobot.Device{buttonred, buttongreen, buttonyellow, buttonblue, buttonwhite, ledred, ledgreen, ledyellow, ledblue, ledwhite},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
