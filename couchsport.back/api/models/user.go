package models

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//User model definition
type User struct {
	Base
	Email       string  `gorm:"unique_index" valid:"email,required" json:"email"`
	Password    string  `valid:"required,length(8|255)" json:"password"`
	PasswordTmp string  `gorm:"-" json:"password_tmp"`
	Profile     Profile `valid:"-" gorm:"foreignkey:ProfileID;constraint:OnDelete:CASCADE;association_autocreate:false;save_associations:false;association_save_reference:true;" json:"profile"`
	ProfileID   uint    `valid:"numeric" json:"profile_id"`
	// // FollowingPages  []*Page `gorm:"many2many:user_page_follower;"`
	// Friends         []*User `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id;"`
	Type           string `valid:"in(ADMIN|USER)" json:"type"`
	New            bool   `gorm:"-" valid:"-" json:"new"`
	ChangePassword bool   `gorm:"-" valid:"-" json:"change_password"`
}

//Validate model
func (user User) Validate(db *gorm.DB) {
	if user.Email == "" {
		db.AddError(errors.New("email is empty"))
		return
	}

	if user.Password == "" {
		db.AddError(errors.New("password is empty"))
		return
	}

	if user.ProfileID == 0 && !user.New && !user.ChangePassword {
		db.AddError(errors.New("profileID invalid"))
	}
}

//BeforeCreate generate the User ID, set Type to USER and hash the password
func (user *User) BeforeCreate(tx *gorm.DB) error {
	profile := Profile{Email: user.Email}
	user.Profile = profile
	user.Type = "USER"
	user.Password = hashAndSalt([]byte(user.Password))
	return nil
}

//AfterCreate empty the password column for security reasons, sets New to true and update Type to ADMIN if ID = 1
func (user *User) AfterCreate(tx *gorm.DB) error {
	user.Password = ""
	user.New = true
	if user.ID == 1 {
		tx.Model(&user).Update("type", "ADMIN")
	}
	return nil
}

//BeforeUpdate set new password if ChangePassword field is true
func (user *User) BeforeUpdate(tx *gorm.DB) error {
	if user.ChangePassword {
		user.Password = hashAndSalt([]byte(user.Password))
	}
	return nil
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
