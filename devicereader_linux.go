package epoller

import (
	"os"
	"syscall"
)

type epollDeviceReader struct {
	epfd int
	fd   int

	handler EventHandler
}

// OpenAndDispatchEvents Blocking call which opens the character device and starts sending events.
func OpenAndDispatchEvents(devicePath string, handler EventHandler, queueDepth int) error {

	edr := epollDeviceReader{handler: handler}

	defer edr.Close()

	if err := edr.Open(devicePath); err != nil {
		return err
	}

	if err := edr.DispatchEvents(queueDepth); err != nil {
		return err
	}
	return nil
}

// NewDeviceReader Returns a new device reader using the supplied event handler
func NewDeviceReader(handler EventHandler) DeviceReader {
	return &epollDeviceReader{handler: handler}
}

// Open the device
func (edr *epollDeviceReader) Open(devicePath string) (err error) {
	// open the device
	edr.fd, err = syscall.Open(devicePath, os.O_RDONLY, 0666)

	if err != nil {
		return
	}
	return
}

// DispatchEvents Blocking routine which sets up epoll and starts sending through events.
func (edr *epollDeviceReader) DispatchEvents(queueDepth int) error {

	var event syscall.EpollEvent
	var events [queueDepth]syscall.EpollEvent

	if err := syscall.SetNonblock(edr.fd, true); err != nil {
		return err
	}

	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return err
	}
	defer syscall.Close(epfd)

	event.Events = syscall.EPOLLIN
	event.Fd = int32(edr.fd)
	if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, edr.fd, &event); err != nil {
		return err
	}

	var nevents int
	for {

		nevents, err = syscall.EpollWait(epfd, events[:], -1)
		if err != nil {
			return err
		}

		for ev := 0; ev < nevents; ev++ {
			// dispatch this to avoid delays in processing
			edr.notifyHandler(int(events[ev].Fd))
		}

	}
	return nil
}

func (edr *epollDeviceReader) notifyHandler(evfd int) {

	var buf [MaxReadSize]byte
	n, _ := syscall.Read(evfd, buf[:])
	if n > 0 {
		edr.handler(buf[0:n], n)
	}
}

// close the device reader and cleanup
func (edr *epollDeviceReader) Close() error {
	syscall.Close(edr.fd)
	return nil
}
