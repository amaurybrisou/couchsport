package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/amaurybrisou/couchsport.back/api/models"
	"github.com/amaurybrisou/couchsport.back/api/stores"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"net/http"
	"strconv"
)

type userHandler struct {
	Store *stores.StoreFactory
}

//All returns all the users
func (me userHandler) All(w http.ResponseWriter, r *http.Request) {

	keys := r.URL.Query()

	users, err := me.Store.UserStore().All(keys)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	json, err := json.Marshal(users)

	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	fmt.Fprint(w, string(json))

}

//Profile returns the connected user profile
func (me userHandler) Profile(userID uint, w http.ResponseWriter, r *http.Request) {
	profile, err := me.Store.UserStore().GetProfile(userID)
	locale := r.Header.Get("Accept-Language")

	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("could not get profile %s", "").Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(profile)

	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("user.could_not_get_profile", locale, nil)).Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json))

}

//SignUp create a user account
func (me userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	locale := r.Header.Get("Accept-Language")

	if r.Body != nil {
		defer r.Body.Close()
	}

	user, err := me.parseBody(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("invalid_request", locale, nil)).Error(), http.StatusBadRequest)
		return
	}

	user, err = me.Store.UserStore().New(user)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("user.could_not_create", locale, map[string]string{"Error": err.Error()})).Error(), http.StatusForbidden)
		return
	}

	json, err := json.Marshal(user)

	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("user.could_not_create", locale, map[string]string{"Error": err.Error()})).Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json))
}

func (me userHandler) ChangePassword(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true
	locale := r.Header.Get("Accept-Language")

	if r.Body != nil {
		defer r.Body.Close()
	}

	user, err := me.parseBody(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("invalid_request", locale, nil)).Error(), http.StatusBadRequest)
		return
	}

	user, err = me.Store.UserStore().ChangePassword(userID, user)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("invalid_request", locale, nil)).Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(user)

	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("internal_error", locale, nil)).Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(json))

}

//Login authenticate the user
func (me userHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	locale := r.Header.Get("Accept-Language")

	if r.Body != nil {
		defer r.Body.Close()
	}

	user, err := me.parseBody(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("invalid_request", locale, nil)).Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := me.Store.UserStore().GetByEmail(user.Email, false)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("invalid_credentials", locale, nil)).Error(), http.StatusUnauthorized)
		return
	}

	if r := comparePasswords(dbUser.Password, []byte(user.Password)); !r {
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("invalid_credentials", locale, nil)).Error(), http.StatusUnauthorized)
		return
	}

	isLogged, err := me.Store.SessionStore().CreateOrRetrieve(dbUser.ID)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("internal_error", locale, nil)).Error(), http.StatusInternalServerError)
		return
	}

	if !isLogged {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("invalid_credentials", locale, nil)).Error(), http.StatusUnauthorized)
		return
	}

	cookie, err := me.Store.SessionStore().CreateCookie()

	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("internal_error", locale, nil)).Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	type res struct {
		Token string `json:"token"`
		Email string `json:"email"`
	}

	responseBody := res{Token: me.Store.SessionStore().GetToken(), Email: dbUser.Email}

	json, err := json.Marshal(responseBody)

	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("internal_error", locale, nil)).Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(json))
}

//IsLogged is a middleware used to know if user is Logged
func (me userHandler) IsLogged(pass func(userID uint, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		session, err := me.Store.SessionStore().GetSession(r)
		if err != nil {
			log.Error(err)
			http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("please_login", locale, nil)).Error(), http.StatusUnauthorized)
			return
		}

		if session.HasExpired() {
			if ok, err := me.Store.SessionStore().Destroy(r); !ok && err != nil {
				log.Error(err)
				http.Error(w, fmt.Errorf("internal error %s", "").Error(), http.StatusInternalServerError)
				return
			}

			if err != nil {
				log.Error(err)
				http.Error(w, fmt.Errorf("internal error %s", "").Error(), http.StatusInternalServerError)
				return
			}

			log.Error(err)
			http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("session_expired", locale, nil)).Error(), http.StatusUnauthorized)
			return
		}

		pass(session.OwnerID, w, r)
	}
}

//Logout log out the user
func (me userHandler) Logout(_ uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true
	locale := r.Header.Get("Accept-Language")
	success, err := me.Store.SessionStore().Destroy(r)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf(me.Store.Localizer().Translate("internal_error", locale, nil)).Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, `{ "Result" : `+strconv.FormatBool(success)+` }`)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}

func (me userHandler) parseBody(body io.Reader) (models.User, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return models.User{}, err
	}

	var u models.User
	err = json.Unmarshal(b, &u)

	if err != nil {
		return models.User{}, err
	}

	return u, nil
}
