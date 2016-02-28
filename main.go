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
	// msg += "\n"
	_, err := SERIAL.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}
}

func Read() string {
	buf := make([]byte, 128)
	n, err := SERIAL.Read(buf)
	// a := fmt.Sprintf("%q", buf[len(buf)-1])
	if err != nil {
		log.Fatal(err)
	}
	//  else {
	// 	for a[len(a)-1] != "\n" {
	// 		fmt.Println(a[len(a)-1])
	// 		n, err = SERIAL.Read(buf)
	// 	}
	// }

	return string(buf[:n])
}

func TestLeds() {
	for _, p := range "12345677654321" {
		led, _ := strconv.Atoi(string(p))
		ToggleLED(led)
		log.Println(Read())
		time.Sleep(200 * time.Millisecond)
	}
}

func TestServos() {
	for a := 0; a < 650; a += 5 {
		SetServo(0, a)
		time.Sleep(40 * time.Millisecond)
	}
	for a := 650; a > 0; a -= 5 {
		SetServo(0, a)
		time.Sleep(40 * time.Millisecond)
	}
}

func TestColor() {
	n := 0
	for {
		SetColorLed(n, 50, 200, 50)
		n++
		if n > 2 {
			n = 0
		}
		SetColorLed(n, 255, 255, 255)
		time.Sleep(80 * time.Millisecond)
	}
}

func SetColorLed(led int, r int, g int, b int) {
	c := fmt.Sprintf("C:%d:%d:%d:%d\n", led, r, g, b)
	Send(c)
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
	log.Println("Start init")
	c := &serial.Config{Name: "/dev/ttyS1", Baud: 115200}
	var err error
	SERIAL, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	ResetLEDs()
	SetAllColor(50, 250, 30)
	SetServo(0, 350)
	// SERIAL.Flush()
}

func SetAllColor(r int, g int, b int) {
	SetColorLed(0, r, g, b)
	SetColorLed(1, r, g, b)
	SetColorLed(2, r, g, b)
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

	// Test := TestServos

	// go func() {
	// 	for {
	// 		TestColor()
	// 	}
	// }()

	if *interactive == true {
		Shell()
	} else {
		for {
			// Test()
			time.Sleep(10 * time.Millisecond)
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
		// line = strings.Replace(" ", ":", line, -1)
		Send(line + "\n")
		log.Println(Read())
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
