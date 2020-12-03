package models

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

//Conversation model definition
type Conversation struct {
	Base
	From     Profile   `gorm:"foreignkey:FromID" json:"from"`
	FromID   uint      `gorm:"association_autoupdate:false;;association_autosave:false;save_associations:false;association_save_reference:false" json:"from_id"`
	To       Profile   `gorm:"foreignkey:ToID" json:"to"`
	ToID     uint      `gorm:"association_autoupdate:false;;association_autosave:false;save_associations:false;association_save_reference:false" json:"to_id"`
	Messages []Message `gorm:"foreignkey:ConversationID;constraint:OnDelete:CASCADE" json:"messages"`
	New      bool      `gorm:"-" json:"new"`
}

//AfterCreate empty the password column for security reasons, sets New to true and update Type to ADMIN if ID = 1
func (Conversation *Conversation) AfterCreate(tx *gorm.DB) error {
	Conversation.New = true
	return nil
}

//Validate model
func (me *Conversation) Validate(db *gorm.DB) {
	if me.FromID < 1 {
		db.AddError(errors.New("invalid FromID"))
	}

	if me.ToID < 1 {
		db.AddError(errors.New("invalid ToID"))
	}

	if me.ToID == me.FromID {
		db.AddError(errors.New("invalid Conversation"))
	}

	if len(me.Messages) < 1 && me.ID < 1 && !me.New {
		db.AddError(errors.New("invalid Messages"))
	}
}

//AddMessage to the expression Messages
func (me *Conversation) AddMessage(fromID, toID uint, text string) Message {
	m := Message{Text: text, Conversation: *me, FromID: fromID, ToID: toID}
	me.Messages = append(me.Messages, m)
	return m
}

//ToJSON converts model to a json string
func (me *Conversation) ToJSON() (string, error) {
	j, err := json.Marshal(me)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
