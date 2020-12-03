package types

import (
	"io"
	"os"
)

//OsFS describes os Filesystem
type OsFS struct{}

func (OsFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
func (OsFS) MkdirAll(path string) error            { return os.MkdirAll(path, 0700) }

func (OsFS) OpenFile(name string) (io.WriteCloser, error) {
	return os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0666)
}

func (OsFS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
