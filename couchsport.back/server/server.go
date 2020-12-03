package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const prefix = "/api"

type Instance struct {
	Db         *gorm.DB
	HTTPServer *http.Server
	router     *http.ServeMux
}

var s *Instance

//NewInstance creates the app server object
func NewInstance() *Instance {
	if s != nil {
		return s
	}
	r := http.NewServeMux()

	db := mustOpenDb()
	// Setup(s.C, s.Db)
	listen := viper.GetString("API_LISTEN")
	port := viper.GetInt("API_PORT")

	log.Println("Listenning on ", listen, port)
	s = &Instance{
		router:     r,
		HTTPServer: &http.Server{Addr: fmt.Sprintf("%s:%d", listen, port), Handler: r},
		Db:         db,
	}

	return s
}

//Start the current Instance
func (s *Instance) Start() {
	go func() {
		if err := s.HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

//Shutdown http server
func (s *Instance) Shutdown() {
	if s.Db != nil {
		// s.Db.Close()
	}

	if s.HTTPServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := s.HTTPServer.Shutdown(ctx)
		if err != nil {
			cancel()
			log.Panic(err)
		} else {
			cancel()
			s.HTTPServer = nil
		}
	}
}

// RegisterHandler registers api handlers
func (s *Instance) RegisterHandler(path string, handler http.HandlerFunc) {
	env := viper.GetString("ENV")
	log.Infof("registering handler at path in %s environment %s, cors is enabled in dev", prefix+path, env)
	if env == "development" {
		handler = s.enableCors(handler)
	}
	s.router.Handle(prefix+path, handler)
}

func (s *Instance) enableCors(pass http.HandlerFunc) http.HandlerFunc {
	port := viper.GetInt("PORT")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:"+strconv.Itoa(port+1))
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:"+strconv.Itoa(port+1))
		w.Header().Set("Access-Control-Request-Headers", "X-Requested-With")
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		pass(w, r)
	}
}
