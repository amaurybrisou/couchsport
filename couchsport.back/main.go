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
	env := viper.GetString("ENV")
	viper.AutomaticEnv()

	languageFiles := viper.GetStringSlice("LANGUAGE_FILES")
	localizer := localizer.NewLocalizer(languageFiles)

	srv := server.NewInstance()

	storeFactory := stores.NewStoreFactory(srv.Db, localizer)
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

	srv.RegisterHandler("/ws", handlerFactory.WsHandler().EntryPoint)

	srv.RegisterHandler("/languages", handlerFactory.LanguageHandler().All)
	srv.RegisterHandler("/activities", handlerFactory.ActivityHandler().All)

	srv.RegisterHandler("/conversations/message/send", handlerFactory.ConversationHandler().HandleMessage)
	srv.RegisterHandler("/conversations/delete", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ConversationHandler().Delete),
	)

	srv.RegisterHandler("/pages", handlerFactory.PageHandler().All)
	srv.RegisterHandler("/pages/new", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().New),
	)
	srv.RegisterHandler("/pages/update", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().Update),
	)
	srv.RegisterHandler("/pages/publish", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().Publish),
	)
	srv.RegisterHandler("/pages/delete", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().Delete),
	)

	srv.RegisterHandler("/images/delete", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ImageHandler().Delete),
	)

	// srv.RegisterHandler("/users", handlerFactory.UserHandler().All)

	srv.RegisterHandler("/profiles/update", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ProfileHandler().Update),
	)
	srv.RegisterHandler("/profiles/mine", handlerFactory.UserHandler().IsLogged(
		handlerFactory.UserHandler().Profile),
	)
	srv.RegisterHandler("/profiles/pages", handlerFactory.UserHandler().IsLogged(
		handlerFactory.PageHandler().ProfilePages),
	)
	srv.RegisterHandler("/profile/conversations", handlerFactory.UserHandler().IsLogged(
		handlerFactory.ConversationHandler().ProfileConversations),
	)

	srv.RegisterHandler("/login", handlerFactory.UserHandler().Login)
	srv.RegisterHandler("/signup", handlerFactory.UserHandler().SignUp)
	srv.RegisterHandler("/logout", handlerFactory.UserHandler().IsLogged(
		handlerFactory.UserHandler().Logout),
	)
	srv.RegisterHandler("/users/change-password", handlerFactory.UserHandler().IsLogged(
		handlerFactory.UserHandler().ChangePassword),
	)

	staticSrv := server.InitStatic()

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

		if err := srv.HTTPServer.Shutdown(ctx); err != nil {
			log.Panic(err)
		}
		log.Info("HTTPServer gracefully closed")

		close(signalDone)
	}()

	go func() {
		if err := staticSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	srv.Start()

	<-signalDone
}
