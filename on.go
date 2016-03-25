package main

import (
	//"encoding/hex"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"strconv"
	//"encoding/binary"

	"time"
)

func main() {
	port := get_port("/dev/ttyACM0", 19200)

	freq := 10

	strip := NewRGBLedStrip(port, freq, 7, 6, 5)
	strip.setRGB("ffffff")
	time.Sleep(time.Duration(2000) * time.Millisecond)
	strip.setRGB("00ff00")
	time.Sleep(time.Duration(2000) * time.Millisecond)
	strip.setRGB("999999")
	time.Sleep(time.Duration(2000) * time.Millisecond)

}

type RGBLedStrip struct {
	port *serial.Port
	ledR Led
	ledG Led
	ledB Led
	freq int
}

func NewRGBLedStrip(port *serial.Port, freq int, pinR int, pinG int, pinB int) *RGBLedStrip {

	s := new(RGBLedStrip)

	s.port = port
	s.freq = freq
	s.ledR = *NewLed(pinR, port, freq)
	s.ledG = *NewLed(pinG, port, freq)
	s.ledB = *NewLed(pinB, port, freq)

	return s
}

func (s *RGBLedStrip) setRGB(hexcode string) {
	//fmt.Printf("send to 1")
	s.ledR.set_value(hexcode[0:2])
	//fmt.Printf("send to 2")
	s.ledG.set_value(hexcode[2:4])
	//fmt.Printf("send to 3")
	s.ledB.set_value(hexcode[4:6])
}

func get_port(name string, baud int) *serial.Port {
	c := &serial.Config{Name: name, Baud: baud}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	return port
}

func send_command(port *serial.Port, command string) string {

	// Write to port
	n, err := port.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}

	//Read back from port
	buf := make([]byte, 128)
	n, err = port.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%q", buf[:n])
}

func (c *PinController) on() string {
	command := "gpio set " + strconv.Itoa(c.pin) + "\r"
	return send_command(c.port, command)
}

func (c *PinController) off() string {
	command := "gpio clear " + strconv.Itoa(c.pin) + "\r"
	return send_command(c.port, command)
}

// LED

type Led struct {
	pin        int
	channel    chan string
	controller PinController
	freq       int
}

func NewLed(pin int, port *serial.Port, freq int) *Led {
	l := new(Led)
	l.freq = freq

	l.pin = pin
	l.channel = make(chan string)
	l.controller = *NewPinController(pin, port)

	go l.start()
	return l
}


func (l *Led) start() {

	s := ""
	num := 0
	on_time := 0
	off_time := 0

	pulse := 1.0 /float64(l.freq)

	pulse = 0.01
	fmt.Printf("%f\n", pulse)
	for {

		select {
		case s = <-l.channel:

			num = hexToInt(s)

			on_ratio := (float64(num) / 255.0)
			on_time := float64(pulse) * on_ratio
			off_time := float64(pulse) - on_time

			fmt.Println("Setting Led ", l.pin, "to ", num, "onetime :", on_time, "offtime :", off_time)
		default:

		}




		//t0 := time.Now()


		l.controller.on()
		time.Sleep(time.Duration(on_time) * time.Millisecond)
		l.controller.off()
		time.Sleep(time.Duration(off_time) * time.Millisecond)

                    //t1 := time.Now()
                    //fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
	}
}

func hexToInt(in string) int {
	c, err := strconv.ParseInt(in, 16, 32)
	if err != nil {

		c = 0 // OR SOMETHING BETTER!
	}
	return int(c)
}

func (l *Led) set_value(value string) {
	l.channel <- value
}


type PinController struct {
	pin  int
	port *serial.Port
}

func NewPinController(pin int, port *serial.Port) *PinController {
	p := new(PinController)
	p.pin = pin
	p.port = port
	return p
}

