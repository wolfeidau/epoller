//
// Package epoller uses the Linux epoll syscall to read data from character devices.
//
// func lineHandler(buf []byte, n int) {
// 	log.Printf("data % X", buf)
// }
//
// func main() {
//
// 	if err := OpenAndDispatchEvents("/dev/kmsg", lineHandler); err != nil {
// 		log.Fatalf("Error opening device reader %v", err)
// 	}
//
// }
//
package epoller

const (
	// MaxReadSize maximum read size
	MaxReadSize = 1024
)

// EventHandler is used to subscribe to handle event data.
type EventHandler func(slice []byte, n int)

// DeviceReader is a simple character device reader using epoll.
type DeviceReader interface {
	Open(devicePath string) (err error)
	DispatchEvents(queueDepth int) (err error)
	Close() error
}
