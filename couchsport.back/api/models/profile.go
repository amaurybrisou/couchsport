package models

import (
	"errors"

	"gorm.io/gorm"
)

//Profile definition
type Profile struct {
	Base
	Username     string `valid:"name" gorm:"type:varchar(50);" json:"username"`
	Country      string `valid:"name" gorm:"type:varchar(50);" json:"country"`
	Firstname    string `valid:"name" gorm:"type:varchar(50);" json:"firstname"`
	Lastname     string `valid:"name" gorm:"type:varchar(50);" json:"lastname"`
	Email        string `valid:"email" json:"email"`
	StreetNumber uint   `valid:"numeric" json:"street_number"`
	StreetName   string `valid:"text" json:"street_name"`
	City         string `valid:"name" gorm:"type:varchar(50);" json:"city"`
	Gender       string `valid:"in(Male|Female)" json:"gender"`
	Phone        string `valid:"alphanum" json:"phone"`
	ZipCode      string `valid:"zipcode" json:"zip_code"`
	Avatar       string `valid:"requri" json:"avatar"`
	AvatarFile   string `gorm:"-" valid:"-" json:"avatar_file"`
	New          bool   `gorm:"-" json:"new"`
	// User                                                                             User
	// OwnerID                                                                          uint        `gorm:"association_autoupdate:false;association_autocreate:false"`
	OwnedPages []Page `gorm:"foreignkey:OwnerID;association_autoupdate:false;association_autocreate:false" json:"owned_pages"`

	Activities []*Activity `gorm:"many2many:profile_activities;association_autoupdate:false;association_autocreate:false" json:"activities"`
	Languages  []*Language `gorm:"many2many:profile_languages;association_autoupdate:false;association_autocreate:false" json:"languages"`
}

//AfterCreate sets New to true
func (p *Profile) AfterCreate(tx *gorm.DB) error {
	p.New = true
	return nil
}

//Validate model
func (p *Profile) Validate(db *gorm.DB) {
	if len(p.Username) > 255 {
		db.AddError(errors.New("invalid Username"))
		return
	}

	if len(p.Firstname) > 255 {
		db.AddError(errors.New("invalid Firstname"))
		return
	}

	if len(p.Lastname) > 255 {
		db.AddError(errors.New("invalid Lastname"))
		return
	}

	if len(p.Phone) > 255 {
		db.AddError(errors.New("invalid Phone"))
		return
	}

	if len(p.ZipCode) > 255 {
		db.AddError(errors.New("invalid ZipCode"))
		return
	}

	if len(p.StreetName) > 255 {
		db.AddError(errors.New("invalid StreetName"))
		return
	}

	if len(p.Country) > 255 {
		db.AddError(errors.New("invalid Country"))
		return
	}

	if len(p.City) > 255 {
		db.AddError(errors.New("invalid City"))
		return
	}
}
