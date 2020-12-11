package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/amaurybrisou/couchsport.back/api/api_errors"
	"github.com/amaurybrisou/couchsport.back/api/models"
	"github.com/amaurybrisou/couchsport.back/api/stores"
	log "github.com/sirupsen/logrus"
)

type profileHandler struct {
	Store *stores.StoreFactory
}

//Update the user profile
func (me profileHandler) Update(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	profile, err := me.parseBody(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, api_errors.ErrInvalidData.Error(), http.StatusUnprocessableEntity)
		return
	}

	owns, err := me.Store.UserStore().OwnProfile(userID, profile.ID)
	if err != nil {
		log.Error(err)
		http.Error(w, api_errors.ErrInvalidData.Error(), http.StatusUnprocessableEntity)
		return
	}

	if !owns {
		log.Error(err)
		http.Error(w, api_errors.ErrDoesNotOwn.Error(), http.StatusForbidden)
		return
	}

	profile, err = me.Store.ProfileStore().Update(userID, profile)
	if err != nil {
		log.Error(err)
		http.Error(w, api_errors.ErrInvalidData.Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(profile)

	if err != nil {
		log.Error(err)
		http.Error(w, api_errors.ErrInternalError.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, string(json))

}

func (me profileHandler) parseBody(body io.Reader) (models.Profile, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return models.Profile{}, err
	}

	var obj models.Profile
	err = json.Unmarshal(b, &obj)

	if err != nil {
		return models.Profile{}, err
	}

	return obj, nil
}
