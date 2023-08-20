// Package iouringfile provides a fast implementation of os.File and associated functions
// using io_uring if you are on Linux.
// It is a drop-in replacement for the os package's os.File type and associated methods.
// On non-Linux packages, it is a wrapper around the os package.
package iouringfile

import (
	"io"
	"io/fs"
	"os"
)

// File represents an open file descriptor. It implements interfaces
// io.Reader, io.ReaderAt, io.Writer, io.WriterAt, io.Seeker, io.Closer,
// io.Syncer, fs.File and fs.FileInfo.
// This version differs from os.File in that it uses io_uring instead of epoll on Linux.
type File struct {
	f *os.File
}

func (f *File) Chdir() error {
	return f.f.Chdir()
}

func (f *File) Chmod(mode fs.FileMode) error {
	return f.f.Chmod(mode)
}

func (f *File) Chown(uid, gid int) error {
	return f.f.Chown(uid, gid)
}

func (f *File) Close() error {
	return f.f.Close()
}

func (f *File) Fd() uintptr {
	return f.f.Fd()
}

func (f *File) Name() string {
	return f.f.Name()
}

/*
ReadDir reads the contents of the directory associated with the file f and returns a slice of DirEntry values in directory order. Subsequent calls on the same file will yield later DirEntry records in the directory.

If n > 0, ReadDir returns at most n DirEntry records. In this case, if ReadDir returns an empty slice, it will return an error explaining why. At the end of a directory, the error is io.EOF.

If n <= 0, ReadDir returns all the DirEntry records remaining in the directory. When it succeeds, it returns a nil error (not io.EOF).
*/
func (f *File) ReadDir(n int) ([]fs.DirEntry, error) {
	return f.f.ReadDir(n)
}

// ReadFrom implements io.ReaderFrom. It does not use io_uring.
func (f *File) ReadFrom(r io.Reader) (n int64, err error) {
	return f.f.ReadFrom(r)
}

/*
Readdir reads the contents of the directory associated with file and returns a slice of up to n FileInfo values, as would be returned by Lstat, in directory order. Subsequent calls on the same file will yield further FileInfos.

If n > 0, Readdir returns at most n FileInfo structures. In this case, if Readdir returns an empty slice, it will return a non-nil error explaining why. At the end of a directory, the error is io.EOF.

If n <= 0, Readdir returns all the FileInfo from the directory in a single slice. In this case, if Readdir succeeds (reads all the way to the end of the directory), it returns the slice and a nil error. If it encounters an error before the end of the directory, Readdir returns the FileInfo read until that point and a non-nil error.

Most clients are better served by the more efficient ReadDir method.
*/
func (f *File) Readdir(n int) ([]fs.FileInfo, error) {
	return f.f.Readdir(n)
}

/*
Readdirnames reads the contents of the directory associated with file and returns a slice of up to n names of files in the directory, in directory order. Subsequent calls on the same file will yield further names.

If n > 0, Readdirnames returns at most n names. In this case, if Readdirnames returns an empty slice, it will return a non-nil error explaining why. At the end of a directory, the error is io.EOF.

If n <= 0, Readdirnames returns all the names from the directory in a single slice. In this case, if Readdirnames succeeds (reads all the way to the end of the directory), it returns the slice and a nil error. If it encounters an error before the end of the directory, Readdirnames returns the names read until that point and a non-nil error.
*/
func (f *File) Readdirnames(n int) ([]string, error) {
	return f.f.Readdirnames(n)
}

/*
Seek sets the offset for the next Read or Write on file to offset,
interpreted according to whence: 0 means relative to the origin of the file,
1 means relative to the current offset, and 2 means relative to the end.
It returns the new offset and an error, if any.
The behavior of Seek on a file opened with O_APPEND is not specified.
*/
func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	return f.f.Seek(offset, whence)
}

func (f *File) SetDeadline(t int64) error {
	return nil
}

func (f *File) SetReadDeadline(t int64) error {
	return nil
}

func (f *File) SetWriteDeadline(t int64) error {

	return nil
}

func (f *File) Stat() (fs.FileInfo, error) {
	return f.f.Stat()
}

func (f *File) Sync() error {
	return f.f.Sync()
}

func (f *File) Truncate(size int64) error {
	return f.f.Truncate(size)
}

// Create creates the named file with mode 0666 (before umask), truncating it if it already exists.
func Create(name string) (*File, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &File{f: f}, nil
}

// Open opens the named file for reading. If successful, methods on the returned
// file can be used for reading; the associated file descriptor has mode O_RDONLY.
func Open(name string) (*File, error) {
	f, err := os.OpenFile(name, 0, 0)
	if err != nil {
		return nil, err
	}

	return &File{f: f}, nil
}

// OpenFile is the generalized open call; most users will use Open or Create instead.
func OpenFile(name string, flag int, perm fs.FileMode) (*File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return &File{f: f}, nil
}
