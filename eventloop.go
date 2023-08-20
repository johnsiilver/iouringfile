//go:build linux

package iouringfile

import (
	"fmt"
	"os"
	"runtime"

	"github.com/iceber/iouring-go"
)

// eventType is the type of event that is sent to the event loop.
type eventType int

const (
	etUnknown eventType = iota
	// etRead is sent when we want to read from a file.
	etRead
	// etWrite is sent when we want to write to a file.
	etWrite
	// etClose is sent when we want to close a file.
	etClose
)

type event struct {
	eventType eventType

	rwEvent rwEvent

	result chan iouring.Result
}

type rwEvent struct {
	fd     uintptr
	offset int64
	b      []byte
}

func submitRWEvent(eventType eventType, f *os.File, b []byte, offset int64) event {
	e := event{
		eventType: eventType,
		rwEvent: rwEvent{
			fd:     f.Fd(),
			offset: offset,
			b:      b,
		},
		result: make(chan iouring.Result, 1),
	}
	eventCh <- e
	return e
}

var eventCh = make(chan event, 100)

func init() {
	go eventLoop()
}

func eventLoop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	iour, err := iouring.New(100, iouring.WithSQPoll())
	if err != nil {
		panic(fmt.Sprintf("new IOURing error: %v", err))
	}
	defer iour.Close()

	for {
		e := <-eventCh
		switch e.eventType {
		case etRead:
			if e.rwEvent.offset == -1 {
				req := iouring.Read(e.rwEvent.fd, e.rwEvent.b)
				if _, err := iouring.SubmitRequest(req, e.result); err != nil {
					panic(err)
				}
			} else {
				req := iouring.Pread(e.rwEvent.fd, e.rwEvent.b, e.rwEvent.offset)
				if _, err := iouring.SubmitRequest(req, e.result); err != nil {
					panic(err)
				}
			}
		case etWrite:
			if e.rwEvent.offset == -1 {
				req := iouring.Write(e.rwEvent.fd, e.rwEvent.b)
				if _, err := iouring.SubmitRequest(req, e.result); err != nil {
					panic(err)
				}
			} else {
				req := iouring.Pwrite(e.rwEvent.fd, e.rwEvent.b, e.rwEvent.offset)
				if _, err := iouring.SubmitRequest(req, e.result); err != nil {
					panic(err)
				}
			}
		}
	}
}
