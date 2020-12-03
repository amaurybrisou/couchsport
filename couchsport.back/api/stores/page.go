package stores

import (
	"net/url"
	"strconv"

	"github.com/amaurybrisou/couchsport.back/api/models"
	"github.com/amaurybrisou/couchsport.back/api/utils"
	"gorm.io/gorm"
)

type pageStore struct {
	Db           *gorm.DB
	FileStore    fileStore
	ImageStore   imageStore
	ProfileStore profileStore
}

//Migrate creates the model schema in database
func (me pageStore) Migrate() {
	err := me.Db.AutoMigrate(&models.Page{})
	if err != nil {
		panic(err)
	}
	// me.Db.Model(&models.Page{}).AddForeignKey("owner_id", "profiles(id)", "NO ACTION", "RESTRICT")
}

//All returns all pages in Database
//Additional keys (url.Values) can be specified :
//followers : returns pages followers
//profile : returns pages profiles
//id: fetch a specific page
func (me pageStore) All(keys url.Values) ([]models.Page, error) {
	var req = me.Db

	req = req.Preload("Images").Preload("Activities")

	for i, v := range keys {
		switch i {
		case "followers":
			req = req.Preload("Followers")
		case "profile":
			req = req.Preload("Owner").Preload("Owner.Languages")
		case "id":
			req = req.Where("ID= ?", v)
		case "name":
			if v[0] == "random" {
				req = req.Order("RAND() LIMIT 1")
				break
			}
			req = req.Where("Name= ?", v)
		case "owner_id":
			req = req.Where("owner_id = ?", v)
		}
	}

	var pages []models.Page
	if err := req.Find(&pages).Error; err != nil {
		return []models.Page{}, err
	}
	return pages, nil
}

//GetPagesByOwnerID return all profile details
func (me pageStore) GetPagesByOwnerID(profileID uint) ([]models.Page, error) {
	var pages []models.Page
	if err := me.Db.Model(&models.Page{}).Preload("Activities").Preload("Images").Where("owner_id = ?", profileID).Find(&pages).Error; err != nil {
		return nil, err
	}

	return pages, nil
}

//New creates a page
func (me pageStore) New(profileID uint, page models.Page) (models.Page, error) {
	page.New = true

	page.OwnerID = profileID

	directory := "page-" + strconv.FormatUint(uint64(profileID), 10)
	images, err := me.downloadImages(directory, page.Images)
	if err != nil {
		return models.Page{}, err
	}

	page.Images = images

	if err := me.Db.Create(&page).Error; err != nil {
		return models.Page{}, err
	}

	return page, nil
}

//Update the page
func (me pageStore) Update(userID uint, page models.Page) (models.Page, error) {
	page.New = false

	if len(page.Images) > 0 {
		directory := "page-" + strconv.FormatUint(uint64(userID), 10)
		images, err := me.downloadImages(directory, page.Images)
		if err != nil {
			return models.Page{}, err
		}

		page.Images = images
		// me.Db.Model(&page).Association("Images").Replace(images)
	}

	me.Db.Unscoped().Table("page_activities").Where("activity_id NOT IN (?)", me.getActivitiesIDS(page.Activities)).Where("page_id = ?", page.ID).Delete(&models.Image{})

	if err := me.Db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", page.ID).Updates(&page).Error; err != nil {
		return models.Page{}, err
	}

	return page, nil
}

//Delete set page.DeletedAt to time.Now() // soft delete thus
func (me pageStore) Delete(userID, pageID uint) (bool, error) {
	if err := me.Db.Exec("DELETE FROM pages WHERE id = ?", pageID).Error; err != nil {
		return false, err
	}

	return true, nil
}

//Publish set page.Public field to 0 or 1
func (me pageStore) Publish(userID, pageID uint, status bool) (bool, error) {
	if err := me.Db.Table("pages").Where("id = ?", pageID).Update("Public", status).Error; err != nil {
		return false, err
	}

	return true, nil
}

// func (me pageStore) getImagesIDS(images []models.Image) []uint {
// 	tmp := []uint{0}
// 	for _, el := range images {
// 		tmp = append(tmp, el.ID)
// 	}
// 	return tmp
// }

func (me pageStore) getActivitiesIDS(activities []*models.Activity) []uint {
	tmp := []uint{0}
	for _, l := range activities {
		tmp = append(tmp, (*l).ID)
	}
	return tmp
}

func (me pageStore) downloadImages(directory string, images []models.Image) ([]models.Image, error) {
	var tmpImages []models.Image
	if len(images) > 0 {
		for idx, i := range images {
			if !i.IsValid() {
				continue
			}
			if i.File != "" && idx < 6 {

				//decode b64 string to bytes
				mime, buf, err := utils.B64ToImage(i.URL)
				if err != nil {
					continue
				}

				img, err := utils.ImageToTypedImage(mime, buf)
				if err != nil {
					continue
				}

				i.File, err = utils.Sanitize(i.File)
				if err != nil {
					i.File = utils.RandStringBytesMaskImprSrc(len(i.File)) + "." + mime
				}

				filename, err := me.FileStore.Save(directory, i.File, img)
				if err != nil {
					continue
				}

				i.File = ""
				i.URL = filename

				tmpImages = append(tmpImages, i)
			} else {
				tmpImages = append(tmpImages, i)
			}
		}
	}
	return tmpImages, nil
}
