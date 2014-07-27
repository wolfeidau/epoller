# epoller

This is just a simple reader which uses epoll on Linux, this is currently being used to read from simple character devices in linux. It also intended as a experiment for those interested in messing around with syscalls in golang.

# why?

Why can't we just use the normal file reader in golang?

Because some very basic character devices don't implement all the neccessary calls to support it.

# Usage

```go
func lineHandler(buf []byte, n int) {
	log.Printf("data % X", buf)
}

func main() {

	if err := epoller.OpenAndDispatchEvents("/dev/kmsg", lineHandler); err != nil {
		log.Fatalf("Error opening device reader %v", err)
	}

}
```

# Licence

Copyright (c) 2014 Mark Wolfe
Licensed under the MIT license.