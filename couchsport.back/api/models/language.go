package models

//Language model definition
type Language struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	Name       string     `gorm:"unique_index;" json:"name"`
	NativeName string     `gorm:"unique_index;" json:"native_name"`
	Profiles   []*Profile `gorm:"many2many:profile_languages;" json:"-"`
}
