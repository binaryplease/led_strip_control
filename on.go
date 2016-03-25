package main

import (
	"fmt"
	"log"

	"github.com/tarm/serial"
	"strconv"
)

func main() {
	port := get_port("/dev/ttyACM0", 19200)
	set_on(7, port)
}

func get_port(name string, baud int) *serial.Port {
	c := &serial.Config{Name: name, Baud: baud}
	port, err := serial.OpenPort(c)
	if err != nil { log.Fatal(err) }
	return port
}

func send_command(port *serial.Port, command string) string {

	// Write to port
	n, err := port.Write([]byte(command))
	if err != nil { log.Fatal(err) }

	//Read back from port
	buf := make([]byte, 128)
	n, err = port.Read(buf)
	if err != nil { log.Fatal(err) }
	return fmt.Sprintf("%q", buf[:n])
}

func set_on(pin int, port *serial.Port) string {
	command := "gpio set " + strconv.Itoa(pin) + "\r"
	return send_command(port, command)
}

func set_off(pin int, port *serial.Port) string {
	command := "gpio clear" + strconv.Itoa(pin) + "\r"
	return send_command(port, command)
}
