package stores

import (
	"fmt"
	"net/http"
	"time"

	"github.com/amaurybrisou/couchsport.back/api/models"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const tokenKey = "user-token"
const sessionValidity = 60 * 60

type sessionStore struct {
	Db     *gorm.DB
	token  string
	userID uint
}

func (me sessionStore) Migrate() {
	err := me.Db.AutoMigrate(&models.Session{})
	if err != nil {
		panic(err)
	}
	// me.Db.Model(&models.Session{}).AddForeignKey("owner_id", "users(id)", "CASCADE", "CASCADE")

}

func (me *sessionStore) CreateOrRetrieve(userID uint) (bool, error) {
	me.userID = userID

	out := models.Session{}
	if err := me.Db.Where("owner_id = ?", userID).Where("(UNIX_TIMESTAMP(expires) - UNIX_TIMESTAMP()) > 0").First(&out).Error; err == gorm.ErrRecordNotFound {
		// record not found => remove all from user and create fresh session
		ok, err := me.DestroyAllByUserID(userID)
		if err != nil {
			return ok, err
		}

		token, err := uuid.NewV4()
		if err != nil {
			return false, err
		}

		session := models.Session{
			SessionID: token.String(),
			OwnerID:   userID,
			Expires:   time.Now().Add(time.Duration(sessionValidity) * time.Second),
			Validity:  sessionValidity,
		}

		if err := me.Db.Create(&session).Error; err != nil {
			return false, err
		}

		me.token = session.SessionID
		return true, nil
	}

	me.token = out.SessionID

	return true, nil
}

func (me *sessionStore) GetSession(r *http.Request) (*models.Session, error) {

	cookie, err := me.GetCookieFromRequest(r)
	if err != nil {
		return nil, err
	}

	if cookie.Value == "" {
		return nil, http.ErrNoCookie
	}

	var session = models.Session{}
	if errs := me.Db.Where("session_id = ?", cookie.Value).First(&session).Error; err != nil {
		log.Errorln(err)
		return nil, errs
	}

	me.token = session.SessionID
	me.userID = session.OwnerID

	return &session, nil
}

func (me *sessionStore) GetCookieFromRequest(r *http.Request) (*http.Cookie, error) {

	c, err := r.Cookie(tokenKey)
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status

			return nil, http.ErrNoCookie
		}
		// For any other type of error, return a bad request status
		return nil, http.ErrNoCookie
	}

	return c, nil

}

func (me *sessionStore) Destroy(r *http.Request) (bool, error) {
	if me.token == "" {
		return false, http.ErrNoCookie
	}

	if err := me.Db.Where("owner_id = ?", me.userID).Delete(&models.Session{}).Error; err != nil {
		log.Errorln(err)
		return false, err
	}

	return true, nil
}

func (me *sessionStore) DestroyAllByUserID(userID uint) (bool, error) {
	if err := me.Db.Where("owner_id = ?", userID).Delete(&models.Session{}).Error; err != nil {
		log.Errorln(err)
		return false, err
	}

	return true, nil
}

func (me *sessionStore) CreateCookie() (*http.Cookie, error) {
	if me.token == "" {
		return nil, fmt.Errorf("cannot generate cookie without token")
	}

	return &http.Cookie{
		Name:    tokenKey,
		Value:   me.token,
		Expires: time.Now().Add(sessionValidity * time.Second),
	}, nil
}

func (me sessionStore) GetToken() string {
	return me.token
}
