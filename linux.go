//go:build linux

package iouringfile

import (
	"io/fs"
	"os"
	"unsafe"
)

/*
Read reads up to len(b) bytes from the File and stores them in b. It returns the number of bytes read and any error encountered. At end of file, Read returns 0, io.EOF.
*/
func (f *File) Read(b []byte) (n int, err error) {
	e := submitRWEvent(etRead, f.f, b, -1)
	i, err := <-e.result.ReturnInt()
	return i, err
}

/*
ReadAt reads len(b) bytes from the File starting at byte offset off. It returns the number of bytes read and the error, if any. ReadAt always returns a non-nil error when n < len(b). At end of file, that error is io.EOF.
*/
func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	e := submitRWEvent(etRead, f.f, b, off)
	i, err := <-e.result.ReturnInt()
	return i, err
}

func (f *File) Write(b []byte) (n int, err error) {
	e := submitRWEvent(etWrite, f.f, b, -1)
	i, err := <-e.result.ReturnInt()
	return i, err
}

// WriteAt writes len(b) bytes to the File starting at byte offset off.
// It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(b).
//
// If file was opened with the O_APPEND flag, WriteAt returns an error.
func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	e := submitRWEvent(etWrite, f.f, b, off)
	i, err := <-e.result.ReturnInt()
	return i, err
}

// WriteString is like Write, but writes the contents of string s rather than
// a slice of bytes.
func (f *File) WriteString(s string) (n int, err error) {
	b := unsafe.Slice(unsafe.StringData(s), len(s))

	e := submitRWEvent(etWrite, f.f, b, -1)
	i, err := <-e.result.ReturnInt()
	return i, err
}

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error.
// This version differs from os.Readfile in that it uses io_uring.
func ReadFile(name string) ([]byte, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}

	b := make([]byte, fi.Size())

	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	e := submitRWEvent(etRead, f, b, -1)
	_, err := <-e.result.ReturnInt()
	return b, err
}

func WriteFile(name string, data []byte, perm fs.FileMode) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	e := submitRWEvent(etWrite, f, data, -1)
	_, err := <-e.result.ReturnInt()
	return err
}
