package handlers

import (
	"github.com/amaurybrisou/couchsport.back/api/stores"
	"github.com/amaurybrisou/couchsport.back/localizer"
	"github.com/gorilla/websocket"
)

//HandlerFactory hols all the handler of the application
type HandlerFactory struct {
	wsHandler           wsHandler
	activityHandler     activityHandler
	imageHandler        imageHandler
	languageHandler     languageHandler
	pageHandler         pageHandler
	profileHandler      profileHandler
	userHandler         userHandler
	conversationHandler conversationHandler
	localizer           *localizer.Localizer
}

//NewHandlerFactory generates the handlerFactory holding every handler in the application
func NewHandlerFactory(storeFactory *stores.StoreFactory, localizer *localizer.Localizer, wsUpgrader *websocket.Upgrader) *HandlerFactory {

	return &HandlerFactory{
		localizer:           localizer,
		wsHandler:           wsHandler{WsUpgrader: wsUpgrader, Stores: storeFactory},
		activityHandler:     activityHandler{Stores: storeFactory},
		imageHandler:        imageHandler{Store: storeFactory},
		languageHandler:     languageHandler{Store: storeFactory},
		pageHandler:         pageHandler{Store: storeFactory},
		profileHandler:      profileHandler{Store: storeFactory},
		userHandler:         userHandler{Store: storeFactory},
		conversationHandler: conversationHandler{Store: storeFactory},
	}
}

//Localizer returns the applicatioin Localizer
func (me HandlerFactory) Localizer() *localizer.Localizer {
	return me.localizer
}

//WsHandler returns the applicatioin Upgrader
func (me HandlerFactory) WsHandler() *wsHandler {
	return &me.wsHandler
}

//ActivityHandler returns the applicatioin ActivityHandler
func (me HandlerFactory) ActivityHandler() *activityHandler {
	return &me.activityHandler
}

//ImageHandler returns the applicatioin ImageHandler
func (me HandlerFactory) ImageHandler() *imageHandler {
	return &me.imageHandler
}

//LanguageHandler returns the applicatioin LanguageHandler
func (me HandlerFactory) LanguageHandler() *languageHandler {
	return &me.languageHandler
}

//PageHandler returns the applicatioin PageHandler
func (me HandlerFactory) PageHandler() *pageHandler {
	return &me.pageHandler
}

//ProfileHandler returns the applicatioin ProfileHandler
func (me HandlerFactory) ProfileHandler() *profileHandler {
	return &me.profileHandler
}

//UserHandler returns the applicatioin UserHandler
func (me HandlerFactory) UserHandler() *userHandler {
	return &me.userHandler
}

//ConversationHandler returns the applicatioin ConversationHandler
func (me HandlerFactory) ConversationHandler() *conversationHandler {
	return &me.conversationHandler
}
