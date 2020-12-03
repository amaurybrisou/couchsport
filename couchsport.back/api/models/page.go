package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

//Page model definition
type Page struct {
	Base
	Name            string      `valid:"text" json:"name"`
	Description     string      `valid:"text" json:"description"`
	LongDescription string      `gorm:"size:512;" valid:"text" json:"long_description"`
	Images          []Image     `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE" json:"images"`
	Lat             float64     `valid:"latitude" json:"lat"`
	Lng             float64     `valid:"longitude" json:"lng"`
	CouchNumber     *int        `valid:"numeric" json:"couch_number"`
	Followers       []*User     `gorm:"many2many:user_page_follower" json:"followers"`
	Owner           Profile     `gorm:"foreignKey:OwnerID;association_autoupdate:false;association_autocreate:false" json:"owner"`
	OwnerID         uint        `json:"owner_id"`
	Public          bool        `gorm:"default:1" json:"public"`
	Activities      []*Activity `gorm:"many2many:page_activities;association_autoupdate:false;association_autocreate:false" json:"activities"`
	New             bool        `gorm:"-" json:"new"`
}

//BeforeCreate is a gorm hook
func (page *Page) BeforeCreate(tx *gorm.DB) error {
	page.CreatedAt = time.Now()
	return nil
}

//AfterCreate is a gorm hook
func (page *Page) AfterCreate(tx *gorm.DB) error {
	page.New = true
	return nil
}

//Validate  tells wheter the page is valid
func (page *Page) Validate(db *gorm.DB) {
	if !page.New && page.ID < 1 {
		db.AddError(errors.New("invalid PageID"))
		return
	}

	if page.Name == "" {
		db.AddError(errors.New("Name is empty"))
		return
	}

	if page.Description == "" {
		db.AddError(errors.New("Description is empty"))
		return
	}

	if len(page.LongDescription) > 512 {
		db.AddError(errors.New("invalid LongDescription"))
		return
	}
}
