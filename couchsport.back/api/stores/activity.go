package stores

import (
	"github.com/amaurybrisou/couchsport.back/api/models"
	"gorm.io/gorm"
)

type activityStore struct {
	Db *gorm.DB
}

//Migrate creates the db table
func (me activityStore) Migrate() {
	err := me.Db.AutoMigrate(&models.Activity{})
	if err != nil {
		panic(err)
	}

	activities := []string{"acrosport", "alpinisme", "apnée", "badminton", "basejump", "basketball", "bmx", "canoëkayak", "canyoning", "course", "coursedorientation", "crosse", "cyclisme", "danse", "équitation", "escalade", "football", "golf", "handball", "hiking", "kitesurfing", "marathon", "paddle", "pêche", "rafting", "roller", "skateboard", "skialpin", "skidefond", "skinautique", "skinordique", "snowboard", "surf", "tennis", "tiràlarc", "ulm", "wakeboard", "yoga", "windsurf"}

	for _, a := range activities {
		me.Db.FirstOrCreate(&models.Activity{Name: a}, models.Activity{Name: a})
	}
}

//All Returns all the activities
func (me activityStore) All() ([]models.Activity, error) {
	var activities []models.Activity
	if err := me.Db.Find(&activities).Error; err != nil {
		return []models.Activity{}, err
	}
	return activities, nil
}
