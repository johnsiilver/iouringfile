//go:build linux

package iouringfile

import (
	"fmt"
	"os"
	//"runtime"

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

	resultCh chan iouring.Result
}

type rwEvent struct {
	fd     uintptr
	offset int64
	b      []byte
}

const resultPoolSize = 100
var resultPool = make(chan chan iouring.Result, resultPoolSize)

func init() {
	for i := 0; i < resultPoolSize; i++ {
		resultPool <- make(chan iouring.Result, 1)
	}
}

func submitRWEvent(eventType eventType, f *os.File, b []byte, offset int64) event {
	e := event{
		eventType: eventType,
		rwEvent: rwEvent{
			fd:     f.Fd(),
			offset: offset,
			b:      b,
		},
		resultCh: <-resultPool,
	}
	eventCh <- e
	return e
}

var eventCh = make(chan event, 100)

func init() {
	go eventLoop()
}

func eventLoop() {
	// This causes the reads to take double the time.
	//runtime.LockOSThread()
	//defer runtime.UnlockOSThread()

	// With iouring.WithSQPoll makes this lock up forever.
	// Also, 100 vs 1 makes no difference for this test, which makes sense without concurrency.
	//iour, err := iouring.New(100, iouring.WithSQPoll())
	iour, err := iouring.New(100)
	if err != nil {
		panic(fmt.Sprintf("new IOURing error: %v", err))
	}
	defer iour.Close()

	for {
		e := <-eventCh
		switch e.eventType {
		case etRead:
			if e.rwEvent.offset == -1 {
				req := iouring.Read(int(e.rwEvent.fd), e.rwEvent.b)
				if _, err := iour.SubmitRequest(req, e.resultCh); err != nil {
					panic(err)
				}
			} else {
				req := iouring.Pread(int(e.rwEvent.fd), e.rwEvent.b, uint64(e.rwEvent.offset))
				if _, err := iour.SubmitRequest(req, e.resultCh); err != nil {
					panic(err)
				}
			}
		case etWrite:
			if e.rwEvent.offset == -1 {
				req := iouring.Write(int(e.rwEvent.fd), e.rwEvent.b)
				if _, err := iour.SubmitRequest(req, e.resultCh); err != nil {
					panic(err)
				}
			} else {
				req := iouring.Pwrite(int(e.rwEvent.fd), e.rwEvent.b, uint64(e.rwEvent.offset))
				if _, err := iour.SubmitRequest(req, e.resultCh); err != nil {
					panic(err)
				}
			}
		}
	}
}
