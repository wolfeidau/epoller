package epoller

import (
	"os"
	"syscall"
)

type epollDeviceReader struct {
	event  syscall.EpollEvent
	events [MaxEpollEvents]syscall.EpollEvent

	epfd int
	fd   int

	handler EventHandler
}

// One stop function for opening and starting device events.
func OpenAndDispatchEvents(devicePath string, handler EventHandler) (err error) {

	edr := epollDeviceReader{handler: handler}

	defer edr.Close()

	if err = edr.Open(devicePath); err != nil {
		return
	}

	if err = edr.DispatchEvents(); err != nil {
		return
	}
	return
}

// Create a new device reader using the supplied event handler
func NewDeviceReader(handler EventHandler) DeviceReader {
	return &epollDeviceReader{handler: handler}
}

// Open the device file and register the epoll events
func (edr *epollDeviceReader) Open(devicePath string) (err error) {
	// open the device
	edr.fd, err = syscall.Open(devicePath, os.O_RDONLY, 0666)

	if err != nil {
		return
	}

	if err = syscall.SetNonblock(edr.fd, true); err != nil {
		return
	}

	edr.epfd, err = syscall.EpollCreate1(0)
	if err != nil {
		return
	}

	edr.event.Events = syscall.EPOLLIN
	edr.event.Fd = int32(edr.fd)
	if err = syscall.EpollCtl(edr.epfd, syscall.EPOLL_CTL_ADD, edr.fd, &edr.event); err != nil {
		return
	}
	return
}

// This is a blocking routine which will start the event loop
func (edr *epollDeviceReader) DispatchEvents() (err error) {

	var nevents int
	for {

		nevents, err = syscall.EpollWait(edr.epfd, edr.events[:], -1)
		if err != nil {
			return
		}

		for ev := 0; ev < nevents; ev++ {
			// dispatch this to avoid delays in processing
			edr.notifyHandler(int(edr.events[ev].Fd))
		}

	}
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
	syscall.Close(edr.epfd)
	syscall.Close(edr.fd)
	return nil
}
