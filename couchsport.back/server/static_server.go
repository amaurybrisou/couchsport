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
func NewStaticHTTPInstance() *StaticInstance {
	if staticInstance != nil {
		return staticInstance
	}

	listen := viper.GetString("STATIC_LISTEN")
	port := viper.GetString("STATIC_PORT")
	r := http.NewServeMux()

	path := viper.GetString("PUBLICPATH")
	log.Infof("[StaticServer] serving files at %s:%s/%s", listen, port, http.Dir(path))
	r.Handle("/uploads/", http.StripPrefix(`/uploads/`, http.FileServer(http.Dir(path+"/uploads"))))

	staticInstance = &StaticInstance{
		router:    r,
		webServer: &http.Server{Addr: listen + ":" + port, Handler: r},
	}

	return staticInstance
}

func NewStaticHTTPSInstance() *StaticInstance {
	if staticInstance != nil {
		return staticInstance
	}
	r := http.NewServeMux()

	path := viper.GetString("PUBLICPATH")
	r.Handle("/uploads/", http.StripPrefix(`/uploads/`, http.FileServer(http.Dir(path+"/uploads"))))

	listen := viper.GetString("STATIC_LISTEN")
	port := viper.GetInt("STATIC_PORT")

	log.Println("[StaticServer] Listenning on https ", listen, port)
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", listen, port),
		Handler:      r,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	staticInstance := &StaticInstance{
		router:    r,
		webServer: srv,
	}

	return staticInstance
}

//Start the current Instance
func (me *StaticInstance) Start() {
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
func (me *StaticInstance) Shutdown(ctx context.Context) (err error) {
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
