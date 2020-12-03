package types

import (
	"io"
	"os"
)

//Filesystem is the interface used in the app (fileStore)
type FileSystem interface {
	OpenFile(name string) (io.WriteCloser, error)
	MkdirAll(path string) error
	Stat(name string) (os.FileInfo, error)
	IsNotExist(error) bool
}
