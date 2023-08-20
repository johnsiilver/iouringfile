package iouringfile

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// globalData is used to prevent the compiler from optimizing away the results.
var globalData []byte

func BenchmarkReadFile(b *testing.B) {
	if runtime.GOOS != "linux" {
		b.Skip("Skipping on non-Linux")
	}

	tests := []struct {
		name     string
		fileSize int
		subTest  string
	}{
		{
			name:     "OS 1KiB",
			fileSize: 1024,
			subTest:  "OSReadFile",
		},
		{
			name:     "OS 16KiB",
			fileSize: 16 * 1024,
			subTest:  "OSReadFile",
		},
		{
			name:     "OS 32KiB",
			fileSize: 32 * 1024,
			subTest:  "OSReadFile",
		},
		{
			name:     "OS 64KiB",
			fileSize: 64 * 1024,
			subTest:  "OSReadFile",
		},
		{
			name:     "OS 128KiB",
			fileSize: 128 * 1024,
			subTest:  "OSReadFile",
		},
		{
			name:     "OS 512KiB",
			fileSize: 512 * 1024,
			subTest:  "OSReadFile",
		},
		{
			name:     "OS 1MiB",
			fileSize: 1024 * 1024,
			subTest:  "OSReadFile",
		},
		{
			name:     "OS 10MiB",
			fileSize: 10 * 1024 * 1024,
			subTest:  "OSReadFile",
		},
		{
			name:     "OS 100MiB",
			fileSize: 100 * 1024 * 1024,
			subTest:  "OSReadFile",
		},
		{
			name:     "OS 1GiB",
			fileSize: 1024 * 1024 * 1024,
			subTest:  "OSReadFile",
		},
		{
			name:     "IOURING 1KiB	",
			fileSize: 1024,
			subTest:  "IOURINGReadFile",
		},
		{
			name:     "IOURING 16KiB",
			fileSize: 16 * 1024,
			subTest:  "IOURINGReadFile",
		},
		{
			name:     "IOURING 32KiB",
			fileSize: 32 * 1024,
			subTest:  "IOURINGReadFile",
		},
		{
			name:     "IOURING 64KiB",
			fileSize: 64 * 1024,
			subTest:  "IOURINGReadFile",
		},
		{
			name:     "IOURING 128KiB",
			fileSize: 128 * 1024,
			subTest:  "IOURINGReadFile",
		},
		{
			name:     "IOURING 512KiB",
			fileSize: 512 * 1024,
			subTest:  "IOURINGReadFile",
		},
		{
			name:     "IOURING 1MiB",
			fileSize: 1024 * 1024,
			subTest:  "IOURINGReadFile",
		},
		{
			name:     "IOURING 10MiB",
			fileSize: 10 * 1024 * 1024,
			subTest:  "IOURINGReadFile",
		},
		{
			name:     "IOURING 100MiB",
			fileSize: 100 * 1024 * 1024,
			subTest:  "IOURINGReadFile",
		},
		{
			name:     "IOURING 1GiB",
			fileSize: 1024 * 1024 * 1024,
			subTest:  "IOURINGReadFile",
		},
	}

	p := filepath.Join(os.TempDir(), "iouringfile_testdata")
	err := os.MkdirAll(p, 0755)
	if err != nil {
		panic(err)
	}

	for _, t := range tests {
		data := make([]byte, t.fileSize)
		rand.Read(data)
		if err := os.WriteFile(filepath.Join(p, t.name), data, 0644); err != nil {
			panic(err)
		}
	}
	defer os.RemoveAll(p)

	for _, t := range tests {
		b.Run(t.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if t.subTest == "IOURINGReadFile" {
					var err error
					globalData, err = ReadFile(filepath.Join(p, t.name))
					if err != nil {
						b.Fatal(err)
					}
				} else {
					var err error
					globalData, err = os.ReadFile(filepath.Join(p, t.name))
					if err != nil {
						b.Fatal(err)
					}
				}
			}
		})
	}
}
