package utils

import (
	"fmt"

	"github.com/amaurybrisou/couchsport.back/api/types"
	log "github.com/sirupsen/logrus"
)

//CreateDirIfNotExists create a directory structure if it doesn't exist already
func CreateDirIfNotExists(_os types.FileSystem, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty %s", path)
	}
	if _, err := _os.Stat(path); _os.IsNotExist(err) {
		log.Printf("creating directory %s", path)
		if err := _os.MkdirAll(path); err != nil {
			log.Printf("error creating directory : %s", err)
			return "", err
		}
	} else if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	return path, nil
}
