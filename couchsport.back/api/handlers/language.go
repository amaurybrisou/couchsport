package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amaurybrisou/couchsport.back/api/stores"
	log "github.com/sirupsen/logrus"
)

type languageHandler struct {
	Store *stores.StoreFactory
}

//All returns all languages
func (app languageHandler) All(w http.ResponseWriter, r *http.Request) {
	languages, err := app.Store.LanguageStore().All()
	if err != nil {
		log.Error(err)
		http.Error(w, http.ErrNotSupported.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(languages)

	if err != nil {
		log.Error(err)
		http.Error(w, http.ErrNotSupported.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json))
}
