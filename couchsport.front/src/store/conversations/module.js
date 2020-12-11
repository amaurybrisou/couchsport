import {
  GET_CONVERSATIONS,
  GOT_CONVERSATIONS,
  CONVERSATION_ADD_MESSAGE,
  CONVERSATION_MESSAGE_SENT,
  CONVERSATION_SEND_MESSAGE,
  NEW_CONVERSATION,
  REMOVE_CONVERSATION,
  CONVERSATION_REMOVED,
  MESSAGES_READ
} from 'store/conversations/actions'

import conversationRepo from 'repos/conversation'
import Vue from 'vue'

const state = { status: '', conversations: [], unread: 0 }

const getters = {}

const actions = {
  [GET_CONVERSATIONS]: ({ commit }) => {
    commit(GET_CONVERSATIONS)
    return conversationRepo.mines().then((response) => {
      commit(GOT_CONVERSATIONS, response.data)
      return response
    })
  },
  [CONVERSATION_SEND_MESSAGE]: ({ commit }, message) => {
    commit(CONVERSATION_SEND_MESSAGE)
    return conversationRepo.sendMessage(message).then((response) => {
      commit(CONVERSATION_MESSAGE_SENT, response.data)
      return response
    })
  },
  [REMOVE_CONVERSATION]: ({ commit }, id) => {
    commit(REMOVE_CONVERSATION)
    return conversationRepo.delete(id).then((response) => {
      commit(CONVERSATION_REMOVED, id)
      return response
    })
  }
}

const mutations = {
  [MESSAGES_READ]: (state, conversation_index) => {
    state.unread = 0
    if (conversation_index > -1) {
      Vue.set(state.conversations[conversation_index], 'unread', false)
    }
  },
  [GET_CONVERSATIONS]: (state) => {
    state.status = 'loading'
  },
  [GOT_CONVERSATIONS]: (state, conversations) => {
    state.status = 'get_success'
    state.conversations = conversations
  },
  [CONVERSATION_ADD_MESSAGE]: (state, message) => {
    state.unread++
    message.data = JSON.parse(message.data)
    for (var i in state.conversations) {
      var c = state.conversations[i]
      if (c.id === message.data.conversation_id) {
        Vue.set(state.conversations[i], 'unread', true)
        state.conversations[i].messages.push(message.data)
        break
      }
    }
  },
  [NEW_CONVERSATION]: (state, message) => {
    state.unread++
    let m = JSON.parse(message.data)
    m.unread = true
    state.conversations.push(m)
  },
  [CONVERSATION_SEND_MESSAGE]: (state) => {
    state.status = 'sending'
  },
  [CONVERSATION_MESSAGE_SENT]: (state, message) => {
    for (var i in state.conversations) {
      var c = state.conversations[i]
      if (c.id === message.conversation_id) {
        state.conversations[i].messages.push(message)
        break
      }
    }
    state.status = 'send_success'
  },
  [REMOVE_CONVERSATION]: (state) => {
    state.status = 'removing'
  },
  [CONVERSATION_REMOVED]: (state, id) => {
    // handle ws request to remove
    if (typeof id === 'object') {
      id = parseInt(id.data)
    }
    state.conversations = state.conversations.filter((c) => id !== c.id)
    state.status = 'remove_success'
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
