package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"gopkg.in/readline.v1"

	"github.com/tarm/serial"
)

var SERIAL *serial.Port

func Send(msg string) {
	_, err := SERIAL.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}
}

func Read() string {
	buf := make([]byte, 128)
	n, err := SERIAL.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%q", buf[:n])
}

func TestLeds() {
	for _, p := range "12345677654321" {
		led, _ := strconv.Atoi(string(p))
		ToggleLED(led)
		log.Println(Read())
		time.Sleep(200 * time.Millisecond)
	}
}

func Test() {
	for a := 0; a < 650; a += 5 {
		SetServo(0, a)
		time.Sleep(40 * time.Millisecond)
	}
	for a := 650; a > 0; a -= 5 {
		SetServo(0, a)
		time.Sleep(40 * time.Millisecond)
	}
}

func SetServo(servo int, pos int) {
	c := fmt.Sprintf("S:%d:%d\n", servo, pos)
	Send(c)
}

func ToggleLED(led int) {
	c := fmt.Sprintf("T:%d\n", led)
	Send(c)
}

func ResetLEDs() {
	Send(fmt.Sprintf("R:0\n"))
}

func Init() {
	// c := &serial.Config{Name: "/dev/ttyS1", Baud: 57600}
	c := &serial.Config{Name: "/dev/ttyS1", Baud: 115200}
	var err error
	SERIAL, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	ResetLEDs()
}

func Heartbeat() {
	for {
		time.Sleep(500 * time.Millisecond)
		SetValue(2, 1)
		time.Sleep(30 * time.Millisecond)
		SetValue(2, 0)
		time.Sleep(250 * time.Millisecond)
		SetValue(2, 1)
		time.Sleep(30 * time.Millisecond)
		SetValue(2, 0)
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	SetMode(2, "out")
	SetMode(4, "out")
	SetValue(2, 1)
	Init()
	SetValue(2, 0)

	interactive := flag.Bool("interactive", false, "readline mode")
	flag.Parse()
	log.Println("Ravenor started")

	go Heartbeat()
	go func() {
		for {
			SetValue(4, 0)
			status := Ping()
			if status {
				SetValue(4, 1)
			}
			time.Sleep(15 * time.Second)
		}
	}()

	if *interactive == true {
		Shell()
	} else {
		for {
			Test()
		}
	}
}

func Shell() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	log.SetOutput(rl.Stderr()) // let "log" write to l.Stderr instead of os.Stderr

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		println(line)
	}
}

func ExportPin(pin int) {
	c := fmt.Sprintf("echo '%d' > /sys/class/gpio/export", pin)
	// log.Println(c)
	cmd := exec.Command("sh", "-c", c)
	cmd.Run()
}

func SetMode(pin int, mode string) {
	ExportPin(pin)
	c := fmt.Sprintf("echo '%s' > /sys/class/gpio/gpio%d/direction", mode, pin)
	// log.Println(c)
	cmd := exec.Command("sh", "-c", c)
	cmd.Run()
}

func SetValue(pin int, val int) {
	c := fmt.Sprintf("echo '%d' > /sys/class/gpio/gpio%d/value", val, pin)
	// log.Println(c)
	cmd := exec.Command("sh", "-c", c)
	cmd.Run()
}

func Reset() {
	SetMode(3, "in")
	time.Sleep(100 * time.Millisecond)
	SetMode(3, "out")
}

func Ping() bool {
	url := "http://google.com"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	defer resp.Body.Close()
	return true
}
