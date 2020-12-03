package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/amaurybrisou/couchsport.back/api/models"
	"github.com/amaurybrisou/couchsport.back/api/stores"
	log "github.com/sirupsen/logrus"
)

type imageHandler struct {
	Store *stores.StoreFactory
}

//Delete is called to set DeletedAt field to Now, not deleting the image
func (me imageHandler) Delete(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	image, err := me.parseBody(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	owns, err := me.Store.UserStore().OwnImage(userID, image.OwnerID, image.ID)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	if !owns {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusForbidden)
		return
	}

	result, err := me.Store.ImageStore().Delete(image.ID)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(struct{ Result bool }{Result: result})

	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json))
}

func (me imageHandler) parseBody(body io.Reader) (models.Image, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return models.Image{}, err
	}

	var obj models.Image
	err = json.Unmarshal(b, &obj)

	if err != nil {
		return models.Image{}, err
	}

	return obj, nil
}
