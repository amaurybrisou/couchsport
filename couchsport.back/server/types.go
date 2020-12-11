package server

import (
	"net/http"

	"gorm.io/gorm"
)

const prefix = "/api"

type Instance struct {
	Db        *gorm.DB
	webServer *http.Server
	router    *http.ServeMux
}

var instance *Instance

type StaticInstance struct {
	webServer *http.Server
	router    *http.ServeMux
}

var staticInstance *StaticInstance
