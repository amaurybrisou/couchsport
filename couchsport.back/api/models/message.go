package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

//Message model definition
type Message struct {
	ID             uint         `gorm:"primarykey" json:"id"`
	Email          string       `valid:"email" json:"email"`
	Date           time.Time    `sql:"DEFAULT:NOW()" json:"date"`
	Text           string       `valid:"text,required" json:"text"`
	From           Profile      `gorm:"foreignkey:FromID" json:"from"`
	FromID         uint         `gorm:"required" json:"from_id"`
	To             Profile      `gorm:"foreignkey:ToID" json:"to"`
	ToID           uint         `gorm:"required" json:"to_id"`
	Conversation   Conversation `gorm:"foreignkey:ConversationID;" valid:"numeric,required" json:"conversation"`
	ConversationID uint         `valid:"numeric,required" json:"conversation_id"`
}

//BeforeCreate is a gorm hook
func (message *Message) BeforeCreate(tx *gorm.DB) error {
	message.Date = time.Now()
	return nil
}

//Validate model
func (m *Message) Validate(db *gorm.DB) {
	if m.Text == "" {
		db.AddError(errors.New("Text is empty"))
		return
	}

	if len(m.Text) > 255 {
		db.AddError(errors.New("invalid Text"))
		return
	}

	if m.ConversationID < 1 {
		db.AddError(errors.New("invalid Conversation"))
		return
	}

	if m.FromID < 1 {
		db.AddError(errors.New("invalid FromID"))
		return
	}
}

//SendMessageBodyModel model definition (model used when decoding body for in conversationHandler.HandleMessage)
type SendMessageBodyModel struct {
	Email string `valid:"email,required" json:"email"`
	Text  string `valid:"text,required" json:"text"`
	ToID  uint   `valid:"numeric" json:"to_id"`
}

//Validate use govalidators to check expression values
func (me SendMessageBodyModel) Validate() (bool, error) {
	if me.ToID < 1 {
		return false, fmt.Errorf("invalid ToID: %v", me.ToID)
	}
	return govalidator.ValidateStruct(me)
}
