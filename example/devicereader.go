package main

import (
	"flag"
	"log"

	. "github.com/wolfeidau/epoller"
)

var deviceFlag = flag.String("device", "/dev/gestic", "device to use")

func lineHandler(buf []byte, n int) {
	log.Printf("data % X", buf)
}

func main() {

	if err := OpenAndDispatchEvents(*deviceFlag, lineHandler); err != nil {
		log.Fatalf("Error opening device reader %v", err)
	}

}
