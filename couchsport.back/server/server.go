package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//NewInstance creates the app server object
func NewHTTPInstance() *Instance {
	if instance != nil {
		return instance
	}
	r := http.NewServeMux()

	db := mustOpenDb()

	// Setup(instance.C, instance.Db)

	listen := viper.GetString("API_LISTEN")
	port := viper.GetInt("API_PORT")

	log.Println("Listenning on http ", listen, port)
	instance = &Instance{
		router:    r,
		webServer: &http.Server{Addr: fmt.Sprintf("%s:%d", listen, port), Handler: r},
		Db:        db,
	}

	return instance
}

func NewHTTPSInstance() *Instance {
	if instance != nil {
		return instance
	}
	r := http.NewServeMux()

	db := mustOpenDb()

	// Setup(instance.C, instance.Db)

	listen := viper.GetString("API_LISTEN")
	port := viper.GetInt("API_PORT")

	log.Println("[ApiServer] Listenning on https ", listen, port)
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", listen, port),
		Handler:      r,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	instance := &Instance{
		router:    r,
		webServer: srv,
		Db:        db,
	}

	return instance
}

//Start the current Instance
func (me *Instance) Start() {
	env := viper.GetString("ENV")
	go func() {
		if env != "development" {
			cert := viper.GetString("CERT_PATH")
			key := viper.GetString("CERT_KEY_PATH")

			if err := me.webServer.ListenAndServeTLS(cert, key); err != http.ErrServerClosed {
				log.Fatal(err)
			}
		} else {
			if err := me.webServer.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()
}

//Shutdown http server
func (me *Instance) Shutdown(ctx context.Context) (err error) {
	// if me.Db != nil {
	// 	 me.Db.Close()
	// }

	if me.webServer != nil {
		err = me.webServer.Shutdown(ctx)
		if err != nil {
			return err
		}
		me.webServer = nil
	}

	return nil
}

// RegisterHandler registers api handlers
func (me *Instance) RegisterHandler(path string, handler http.HandlerFunc) {
	env := viper.GetString("ENV")

	corsEnabled := env != "production"
	log.Infof("[ApiServer] Registering handler at path %s in %s, cors: %v", prefix+path, env, corsEnabled)

	if env != "production" {
		handler = enableCors(handler)
	}

	me.router.Handle(prefix+path, httpsWrapper(handler))
}
