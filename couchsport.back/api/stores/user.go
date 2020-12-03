package stores

import (
	"fmt"
	"net/url"

	"github.com/amaurybrisou/couchsport.back/api/models"
	"github.com/amaurybrisou/couchsport.back/api/utils"
	"gorm.io/gorm"
)

type userStore struct {
	Db *gorm.DB
}

func (me userStore) Migrate() {
	err := me.Db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	// me.Db.Model(&models.User{}).AddForeignKey("profile_id", "profiles(id)", "CASCADE", "CASCADE")
}

//All user fetch
func (me userStore) All(keys url.Values) ([]models.User, error) {
	var req = me.Db
	for i, v := range keys {
		switch i {
		case "profile":
			req = req.Preload("Profile")
		case "pages":
			req = req.Preload("OwnedPages")
		case "follow":
			req = req.Preload("FollowingPages")
		case "friends":
			req = req.Preload("Friends")
		case "id":
			req = req.Where("ID= ?", v)
		case "username":
			req = req.Where("username  LIKE ?", v)
		case "email":
			req = req.Where("email LIKE ?", v)
		}
	}

	var users []models.User
	if err := req.Find(&users).Error; err != nil {
		return []models.User{}, err
	}
	return users, nil
}

//New user
func (me userStore) New(user models.User) (models.User, error) {
	user.New = true

	var count int64
	if err := me.Db.Model(&user).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return models.User{}, err
	}

	if count > 0 {
		return models.User{}, fmt.Errorf("user already exist")
	}

	if err := me.Db.Create(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (me userStore) ChangePassword(userID uint, user models.User) (models.User, error) {
	user.ChangePassword = true
	if err := me.Db.Model(&user).Where("id = ?", userID).Update("Password", user.Password).Error; err != nil {
		return user, err
	}
	return user, nil
}

//GetProfile returns the user profile
func (me userStore) GetProfile(userID uint) (models.Profile, error) {
	var out = models.User{}
	if err := me.Db.
		Preload("Profile").
		Preload("Profile.Languages").
		Preload("Profile.Activities").
		Where("id = ?", userID).First(&out).Error; err != nil {
		return out.Profile, err
	}
	return out.Profile, nil
}

func (me userStore) GetByID(userID uint) (models.User, error) {
	var outUser = models.User{}
	if err := me.Db.Model(&models.User{}).Preload("Profile").Where("id = ?", userID).First(&outUser).Error; err != nil {
		return models.User{}, err
	}
	return outUser, nil
}

func (me userStore) GetByEmail(email string, create bool) (models.User, error) {
	var outUser = models.User{}
	if err := me.Db.Where("email = ?", email).First(&outUser).Error; create && err == gorm.ErrRecordNotFound {
		return me.NewWithoutPassword(email)
	} else if err != nil {
		return models.User{}, err
	}
	return outUser, nil
}

//OwnImage tells you wheter the userID owns PageID whose owns the imageID as well
func (me userStore) OwnImage(userID, pageID, imageID uint) (bool, error) {
	if userID < 1 {
		return false, fmt.Errorf("userID cannot be below 1")
	}

	profileID, err := me.GetProfileID(userID)
	if err != nil {
		return false, err
	}

	var count struct{ Count int }
	err = me.Db.Table("images").Select("COUNT(*) AS count").Joins("INNER JOIN pages ON pages.id = images.owner_id").Where("images.id = ? AND pages.id = ? AND pages.owner_id = ?", imageID, pageID, profileID).Find(&count).Error
	if err != nil {
		return false, err
	}

	if count.Count < 1 {
		return false, fmt.Errorf("user %v doesn't own this image %v", userID, imageID)
	}

	return true, nil
}

//OwnPage tells you wheter the userID owns the pageID
func (me userStore) OwnPage(userID, pageID uint) (bool, error) {
	if userID < 1 {
		return false, fmt.Errorf("userID cannot be below 1")
	}

	user, err := me.GetByID(userID)
	if err != nil {
		return false, err
	}

	var count int64
	if err := me.Db.Model(models.Page{}).Where("owner_id = ?", user.ProfileID).Where("id = ?", pageID).Count(&count).Error; err != nil {
		return false, err
	}

	if count < 1 {
		return false, fmt.Errorf("user %v doesn't own this page %v", userID, pageID)
	}

	return true, nil
}

//OwnConversation tells you wheter the userID owns the conversation
func (me userStore) OwnConversation(userID, conversationID uint) (bool, uint, error) {
	if userID < 1 {
		return false, 0, fmt.Errorf("userID cannot be below 1")
	}

	profileID, err := me.GetProfileID(userID)
	if err != nil {
		return false, 0, err
	}

	var count int64
	var conversation models.Conversation
	if err := me.Db.Model(&conversation).
		Select("from_id", "to_id").
		Where("from_id = ? OR to_id = ?", profileID, profileID).
		Where("id = ?", conversationID).
		First(&conversation).
		Count(&count).Error; err != nil {
		return false, 0, err
	}

	if count < 1 {
		return false, 0, fmt.Errorf("user %v isn't part of this conversation %v", userID, conversationID)
	}

	retID := conversation.ToID
	if profileID == conversation.ToID {
		retID = conversation.FromID
	}

	return true, retID, nil
}

//OwnProfile tells you wheter the userID owns the profile
func (me userStore) OwnProfile(userID, profileID uint) (bool, error) {
	if profileID < 1 || userID < 1 {
		return false, fmt.Errorf("profileID or userID cannot be below 1")
	}

	user, err := me.GetByID(userID)
	if err != nil {
		return false, err
	}

	if user.ID != profileID {
		return false, fmt.Errorf("user %v doesn't own this profile %v", profileID, profileID)
	}

	return true, nil
}

//GetProfileID returns the profileID of the submitted userID
func (me userStore) GetProfileID(userID uint) (uint, error) {
	profile, err := me.GetProfile(userID)
	return profile.ID, err
}

// func parseBody(tmp interface{}) (models.User, error) {
// 	fmt.Println(tmp)
// 	r, ok := tmp.(*models.User)

// 	if !ok {
// 		return models.User{}, fmt.Errorf("body is not of type User")
// 	}

// 	return *r, nil
// }

func (me userStore) NewWithoutPassword(email string) (models.User, error) {
	password := utils.RandStringBytesMaskImprSrc(len(email))
	user := models.User{
		Email:       email,
		Password:    password,
		PasswordTmp: password,
	}

	return me.New(user)
}
