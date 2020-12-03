package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

//Session model definition
type Session struct {
	Owner     User      `valid:"-" gorm:"foreign_key:OwnerId;association_autoupdate:false;association_autocreate:false" json:"owner"`
	OwnerID   uint      `valid:"numeric" json:"owner_id"`
	SessionID string    `valid:"uuidv4" json:"session_id"`
	Expires   time.Time `gorm:"default=now" valid:"-" json:"expires"`
	Validity  uint      `valid:"numeric" json:"validity"`
}

//Validate model
func (session *Session) Validate(db *gorm.DB) {
	if session.OwnerID == 0 {
		db.AddError(errors.New("invalid OwnerID"))
		return
	}

	if session.SessionID == "" {
		db.AddError(errors.New("invalid SessionID"))
		return
	}

	if session.Expires.Before(time.Now()) {
		db.AddError(errors.New("sesssion expired"))
		return
	}
}

//HasExpired determines wether the current session has expired or not
func (session *Session) HasExpired() bool {
	if session.Expires.After(time.Now()) {
		return false
	}
	return true
}
