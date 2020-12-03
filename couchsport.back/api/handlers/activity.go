package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amaurybrisou/couchsport.back/api/stores"
	log "github.com/sirupsen/logrus"
)

type activityHandler struct {
	Stores *stores.StoreFactory
}

//All returns all the activities in DB
func (app activityHandler) All(w http.ResponseWriter, r *http.Request) {
	activities, err := app.Stores.ActivityStore().All()
	if err != nil {
		log.Error(err)
		http.Error(w, "error retrieving activities", http.StatusInternalServerError)
	}

	json, err := json.Marshal(activities)
	if err != nil {
		log.Error(err)
		http.Error(w, "error retrieving activities", http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(json))
}
