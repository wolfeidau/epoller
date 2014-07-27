package epoller

import (
	"flag"
	"log"
	"testing"
)

var dev = flag.String("dev", "/dev/kmsg", "device to use")

func TestDeviceReaderOpen(t *testing.T) {

	dr := NewDeviceReader(func(slice []byte, n int) {})

	defer dr.Close()

	if err := dr.Open(*dev); err != nil {
		t.Fatalf("Error opening device reader %v", err)
	}
}

func TestDeviceReaderEvents(t *testing.T) {

	results := make(chan []byte)

	dr := NewDeviceReader(func(slice []byte, n int) {
		results <- slice
	})

	defer dr.Close()

	if err := dr.Open(*dev); err != nil {
		t.Fatalf("Error opening device reader %v", err)
	}

	go readEvents(dr, t)

	data := <-results

	log.Printf("data %x", data)

}

func readEvents(dr DeviceReader, t *testing.T) {
	if err := dr.DispatchEvents(); err != nil {
		t.Fatalf("Error reading events %v", err)
	}
}
