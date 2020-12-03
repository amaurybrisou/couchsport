package models

import (
	"errors"

	"gorm.io/gorm"
)

//Image model definition
type Image struct {
	Base
	URL     string `valid:"requri,required" json:"url"`
	Alt     string `valid:"text" json:"alt"`
	File    string `gorm:"-" json:"file"`
	OwnerID uint   `json:"owner_id"`
}

//IsValid tells is Image is valid
func (image *Image) IsValid() bool {
	if image.URL == "" {
		return false
	}

	if len(image.URL) > 255 && image.File == "" {
		return false
	}

	if len(image.URL) > 255 && image.ID > 0 {
		return false
	}

	if len(image.Alt) > 255 {
		return false
	}

	return true
}

//Validate tells is Image is valid
func (image *Image) Validate(db *gorm.DB) {
	if image.URL == "" {
		db.AddError(errors.New("URL is empty"))
		return
	}

	if len(image.URL) > 255 {
		db.AddError(errors.New("invalid URL"))
		return
	}

	if len(image.Alt) > 255 {
		db.AddError(errors.New("invalid Alt"))
		return
	}
}
