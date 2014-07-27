package epoller

const (
	MaxEpollEvents = 8    // max events to queue
	MaxReadSize    = 1024 // maximum read size
)

type EventHandler func(slice []byte, n int)

type DeviceReader interface {
	Open(devicePath string) (err error)
	DispatchEvents() (err error)
	Close() error
}
