package stores

import (
	"strconv"

	"github.com/amaurybrisou/couchsport.back/api/models"
	"github.com/amaurybrisou/couchsport.back/api/utils"
	"gorm.io/gorm"
)

type profileStore struct {
	Db        *gorm.DB
	FileStore fileStore
}

//Migrate creates the table in database
func (me profileStore) Migrate() {
	err := me.Db.AutoMigrate(&models.Profile{})
	if err != nil {
		panic(err)
	}
}

//GetProfiles returns all profiles in database
func (me profileStore) All() ([]models.Profile, error) {
	var profiles []models.Profile
	if err := me.Db.Model(&models.Profile{}).
		Find(&profiles).Error; err != nil {
		return []models.Profile{}, err
	}
	return profiles, nil
}

//Update the profile
func (me profileStore) Update(userID uint, profile models.Profile) (models.Profile, error) {
	if profile.AvatarFile != "" {

		filename, err := me.saveAvatar(userID, profile.AvatarFile, profile.Avatar)
		if err != nil {
			return models.Profile{}, err
		}
		profile.AvatarFile = ""
		profile.Avatar = filename
	}

	me.Db.Unscoped().
		Table("profile_activities").
		Where("activity_id NOT IN (?)", me.getActivitiesIDS(profile.Activities)).
		Where("profile_id = ?", profile.ID).
		Delete(&models.Profile{})
	me.Db.Unscoped().
		Table("profile_languages").
		Where("language_id NOT IN (?)", me.getLanguagesIDS(profile.Languages)).
		Where("profile_id = ?", profile.ID).
		Delete(&models.Profile{})

	if err := me.Db.
		Model(&profile).Updates(&profile).Error; err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (me profileStore) getLanguagesIDS(languages []*models.Language) []uint {
	tmp := []uint{0}
	for _, l := range languages {
		tmp = append(tmp, (*l).ID)
	}
	return tmp
}

func (me profileStore) getActivitiesIDS(activities []*models.Activity) []uint {
	tmp := []uint{0}
	for _, l := range activities {
		tmp = append(tmp, (*l).ID)
	}
	return tmp
}

func (me profileStore) saveAvatar(profileID uint, filename, b64 string) (string, error) {
	//decode b64 string to bytes
	mime, buf, err := utils.B64ToImage(b64)
	if err != nil {
		return "", err
	}

	img, err := utils.ImageToTypedImage(mime, buf)
	if err != nil {
		return "", err
	}

	filename, err = utils.Sanitize(filename)
	if err != nil {
		filename = utils.RandStringBytesMaskImprSrc(len(filename)) + "." + mime
	}

	directory := "user-" + strconv.FormatUint(uint64(profileID), 10)
	filename, err = me.FileStore.Save(directory, filename, img)
	if err != nil {
		return "", err
	}

	return filename, nil
}

// func (me profileStore) parseBody(tmp interface{}) (models.Profile, error) {
// 	r, ok := tmp.(models.Profile)

// 	if !ok {
// 		return models.Profile{}, fmt.Errorf("body is not of type Profile")
// 	}

// 	return r, nil
// }
