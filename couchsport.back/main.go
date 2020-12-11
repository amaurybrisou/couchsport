package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/viper"

	"github.com/amaurybrisou/couchsport.back/api/handlers"
	"github.com/amaurybrisou/couchsport.back/api/stores"
	"github.com/amaurybrisou/couchsport.back/api/validators"
	"github.com/amaurybrisou/couchsport.back/localizer"
	"github.com/amaurybrisou/couchsport.back/server"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func main() {
	viper.AutomaticEnv()

	env := viper.GetString("ENV")
	languageFiles := viper.GetStringSlice("LANGUAGE_FILES")
	localizer := localizer.NewLocalizer(languageFiles)

	var api *server.Instance
	var staticSrv *server.StaticInstance
	if env != "production" {
		api = server.NewHTTPInstance()
		staticSrv = server.NewStaticHTTPInstance()
	} else {
		api = server.NewHTTPSInstance()
		staticSrv = server.NewStaticHTTPSInstance()
	}

	storeFactory := stores.NewStoreFactory(api.Db, localizer)
	storeFactory.Init(viper.GetBool("POPULATE"))

	allowedDomains := viper.GetStringSlice("ALLOWED_DOMAINS")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header["Origin"]
			log.Println(origin, allowedDomains)
			for _, d := range allowedDomains {
				if strings.Contains(origin[0], d) {
					return true
				}
			}
			return false
		},
	}

	if env == "development" {
		log.Println("enable WebSocket All Origins")
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	}

	handlerFactory := handlers.NewHandlerFactory(storeFactory, localizer, &upgrader)

	validators.Init()

	api.RegisterHandler("/ws", handlerFactory.WsHandler().EntryPoint)

	api.RegisterHandler("/languages", handlerFactory.LanguageHandler().All)
	api.RegisterHandler("/activities", handlerFactory.ActivityHandler().All)

	api.RegisterHandler("/conversations/message/send", handlerFactory.ConversationHandler().HandleMessage)
	api.RegisterHandler("/conversations/delete", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ConversationHandler().Delete),
	)

	api.RegisterHandler("/pages", handlerFactory.PageHandler().All)
	api.RegisterHandler("/pages/new", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().New),
	)
	api.RegisterHandler("/pages/update", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().Update),
	)
	api.RegisterHandler("/pages/publish", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().Publish),
	)
	api.RegisterHandler("/pages/delete", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().Delete),
	)

	api.RegisterHandler("/images/delete", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ImageHandler().Delete),
	)

	// api.RegisterHandler("/users", handlerFactory.UserHandler().All)

	api.RegisterHandler("/profiles/update", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ProfileHandler().Update),
	)
	api.RegisterHandler("/profiles/mine", handlerFactory.UserHandler().IsLogged(
		handlerFactory.UserHandler().Profile),
	)
	api.RegisterHandler("/profiles/pages", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().ProfilePages),
	)
	api.RegisterHandler("/profile/conversations", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ConversationHandler().ProfileConversations),
	)

	api.RegisterHandler("/login", handlerFactory.UserHandler().Login)
	api.RegisterHandler("/signup", handlerFactory.UserHandler().SignUp)
	api.RegisterHandler("/logout", handlerFactory.UserHandler().IsLogged(
		handlerFactory.UserHandler().Logout),
	)
	api.RegisterHandler("/users/change-password", handlerFactory.UserHandler().IsLogged(
		handlerFactory.UserHandler().ChangePassword),
	)

	signalChan := make(chan os.Signal, 1)
	signalDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		<-signalChan
		log.Info("received os.Interrupt signal, stopping services")
		storeFactory.WsStore().Close(signalDone)
		cancel()
		<-signalDone

		if err := staticSrv.Shutdown(ctx); err != nil {
			log.Panic(err)
		}

		log.Info("Static HTTPServer gracefully closed")

		if err := api.Shutdown(ctx); err != nil {
			log.Panic(err)
		}
		log.Info("HTTPServer gracefully closed")

		close(signalDone)
	}()

	staticSrv.Start()
	api.Start()

	<-signalDone
}
