package models

//Activity model definition
type Activity struct {
	ID       uint       `gorm:"primarykey" json:"id"`
	Name     string     `gorm:"unique_index" json:"name"`
	Profiles []*Profile `gorm:"many2many:profile_activities;" json:"profiles"`
	Pages    []*Page    `gorm:"many2many:page_activities;constraint:OnDelete:CASCADE" json:"pages"`
}
