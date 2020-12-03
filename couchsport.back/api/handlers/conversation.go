package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/amaurybrisou/couchsport.back/api/models"
	"github.com/amaurybrisou/couchsport.back/api/stores"
	log "github.com/sirupsen/logrus"
)

type conversationHandler struct {
	Store *stores.StoreFactory
}

func (me conversationHandler) HandleMessage(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	locale := r.Header.Get("Accept-Language")

	if r.Body != nil {
		defer r.Body.Close()
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	sendMessageBody := models.SendMessageBodyModel{}
	err = json.Unmarshal(body, &sendMessageBody)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	_, err = sendMessageBody.Validate()
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	toProfile, err := me.Store.UserStore().GetProfile(sendMessageBody.ToID)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	fromUser, err := me.Store.UserStore().GetByEmail(sendMessageBody.Email, true)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	fromProfile, err := me.Store.UserStore().GetProfile(fromUser.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	if fromProfile.ID == toProfile.ID {
		http.Error(w, fmt.Errorf("%s", "invalid request").Error(), http.StatusBadRequest)
		return
	}

	conversation, err := me.Store.ConversationStore().GetByReferents(fromProfile, toProfile)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	conversation, message, err := me.Store.ConversationStore().AddMessage(conversation, fromProfile.ID, toProfile.ID, sendMessageBody.Email, sendMessageBody.Text)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	message.From = fromProfile

	if fromUser.New {
		go me.Store.MailStore().AccountAutoCreated(fromUser.Email, fromUser.PasswordTmp, locale)
	}

	j, err := json.Marshal(&message)

	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	if !conversation.New {
		me.Store.WsStore().EmitToMutationNamespace(message.ToID, "CONVERSATION_ADD_MESSAGE", string(j), "conversations")
	} else {
		c, err := conversation.ToJSON()
		if err != nil {
			http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
			return
		}
		me.Store.WsStore().EmitToMutationNamespace(message.ToID, "NEW_CONVERSATION", c, "conversations")
	}

	fmt.Fprint(w, string(j))
}

func (me conversationHandler) ProfileConversations(userID uint, w http.ResponseWriter, r *http.Request) {
	profileID, err := me.Store.UserStore().GetProfileID(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	conversations, err := me.Store.ConversationStore().ProfileConversations(profileID)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(&conversations)

	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json))
}

func (me conversationHandler) Delete(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	query := r.URL.Query()

	tmp := query.Get("id")
	if tmp == "" {
		log.Println("id mising")
		http.Error(w, fmt.Errorf("id missing %s", tmp).Error(), http.StatusBadRequest)
		return
	}

	conversationID, err := strconv.Atoi(tmp)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	owns, interlocutorProfileID, err := me.Store.UserStore().OwnConversation(userID, uint(conversationID))
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	if !owns {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusForbidden)
		return
	}

	result, err := me.Store.ConversationStore().Delete(uint(conversationID))
	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	ret, err := json.Marshal(struct{ Result bool }{Result: result})

	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	me.Store.WsStore().EmitToMutationNamespace(interlocutorProfileID, "CONVERSATION_REMOVED", fmt.Sprint(conversationID), "conversations")

	fmt.Fprint(w, string(ret))
}
