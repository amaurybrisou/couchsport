package stores

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/amaurybrisou/couchsport.back/api/types"
	"github.com/amaurybrisou/couchsport.back/api/utils"
	log "github.com/sirupsen/logrus"
)

type fileStore struct {
	PublicPath, ImageBasePath, FilePrefix string
	FileSystem                            types.FileSystem
}

//Save a file on the filesystem at path computed from ImageBasePath + directory + UserID
//directory is prepend before userId
func (me fileStore) Save(directory, filename string, buf io.Reader) (string, error) {

	if filename == "" {
		err := fmt.Errorf("filename is incorrect")
		return "", err
	}

	path := me.ImageBasePath
	if directory != "" {
		path += string(os.PathSeparator) + directory
	}

	fsPath, err := utils.CreateDirIfNotExists(me.FileSystem, filepath.Join(me.PublicPath, path))
	if err != nil {
		return "", err
	}

	filename = me.FilePrefix + filename

	log.Printf("Openning file %s", fsPath+"/"+filename)
	f, err := me.FileSystem.OpenFile(filepath.Join(fsPath, filename))
	if err != nil {
		return "", err
	}

	defer f.Close()
	count, err := io.Copy(f, buf)
	if err != nil {
		return "", err
	}

	if count == 0 {
		return "", fmt.Errorf("file not created, no data to write")
	}

	log.Printf("%d bytes wrote at %s", count, fsPath+"/"+filename)

	return filepath.Join(path, filename), nil
}
