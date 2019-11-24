package steam

import (
	"io"
	"os"
	"time"
	"errors"
)

type File interface {
	io.Reader
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}


func (c DepotChunk) Stat() (info os.FileInfo, err error) {
	return c, nil
}
func (c DepotChunk) IsDir() bool {
	return false
}
func (c DepotChunk) ModTime() time.Time {
	return time.Unix(0, 0)
}
func (c DepotChunk) Mode() os.FileMode {
	return 222
}
func (c DepotChunk) Sys() interface{} {
	return c
}

func (f DepotFile) Stat() (info os.FileInfo, err error) {
	return f, nil
}
func (f DepotFile) Readdir(n int) (chunks []os.FileInfo, err error) {
	if n <= 0 {
		n = len(f.Chunks)
	}
	chunks = make([]os.FileInfo, n)
	for i, c := range f.Chunks {
		if chunks[i], err = c.Stat(); err != nil {
			return
		}
	}
	return
}
func (f DepotFile) IsDir() bool {
	return false;
}
func (f DepotFile) ModTime() time.Time {
	return time.Unix(0, 0)
}
func (f DepotFile) Mode() os.FileMode {
	return 222
}
func (f DepotFile) Sys() interface{} {
	return f
}

func (d Depot) Read(p []byte) (n int, err error) {
	err = errors.New("is a directory")
	return
}
func (d Depot) Readdir(n int) (files []os.FileInfo, err error) {
	files = make([]os.FileInfo, n)
	for i, f := range d.Files {
		if files[i], err = f.Stat(); err != nil {
			return
		}
	}
	return
}
func (d Depot) Stat() (info os.FileInfo, err error) {
	return d, nil
}
func (i Depot) Mode() os.FileMode {
	return os.ModeDir
}
func (i Depot) IsDir() bool {
	return true
}
func (i Depot) Sys() interface{} {
	return i
}
