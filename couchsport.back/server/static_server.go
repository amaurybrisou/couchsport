package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitStatic() (srv http.Server) {
	listen := viper.GetString("STATIC_LISTEN")
	port := viper.GetString("STATIC_PORT")
	r := http.NewServeMux()

	path := viper.GetString("PUBLICPATH")
	logrus.Infof("[StaticServer] serving files at %s:%s/%s", listen, port, http.Dir(path))
	r.Handle("/uploads/", http.StripPrefix(`/uploads/`, http.FileServer(http.Dir(path+"/uploads"))))

	srv = http.Server{Addr: listen + ":" + port, Handler: r}

	return
}
