package server

import (
	"net/http"
	"strconv"

	"github.com/spf13/viper"
)

func httpsWrapper(pass http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		pass(w, r)
	}
}

func enableCors(pass http.HandlerFunc) http.HandlerFunc {
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
