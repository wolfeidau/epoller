package main

import (
	"flag"
	"github.com/wolfeidau/epoller"
	"log"
	"syscall"
)

var deviceFlag = flag.String("device", "/dev/kmsg", "device to use")
var sdoutfd = syscall.Stdout

func lineHandler(buf []byte, n int) {
	syscall.Write(sdoutfd, buf)
}

func main() {

	if err := epoller.OpenAndDispatchEvents(*deviceFlag, lineHandler); err != nil {
		log.Fatalf("Error opening device reader %v", err)
	}

}
