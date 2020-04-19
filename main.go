package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var fake bool
var ledStatus bool

func setLedStatus(s bool) {
	var err2 error
	var f2 *os.File
	var b [1]byte

	if s {
		b[0] = '1'
	} else {
		b[0] = '0'
	}
	if !fake {
		f2, err2 = os.OpenFile("/sys/class/leds/led0/brightness", os.O_WRONLY, 0777)
		if err2 != nil {
			log.Fatal(err2)
		}
		if _, err2 = f2.Write([]byte("1")); err2 != nil {
			f2.Close()
			log.Fatal(err2)
		}
	}
	log.Printf("LED: %t\n", s)
	ledStatus = s
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Switching to: %t\n", !ledStatus)
	setLedStatus(!ledStatus)
}

func main() {
	var err error
	var f1 *os.File

	fakePtr := flag.Bool("fake", false, "Do not operate leds, only log actions.")
	flag.Parse()
	fake = *fakePtr

	log.Printf("Fake mode: %t\n", fake)

	if !fake {
		f1, err = os.OpenFile("/sys/class/leds/led0/trigger", os.O_WRONLY, 0777)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f1.Write([]byte("none")); err != nil {
			f1.Close()
			log.Fatal(err)
		}
	}
	log.Println("Trigger set to 'none'")
	setLedStatus(false)

	http.HandleFunc("/switch", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
