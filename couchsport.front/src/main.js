import Vue from 'vue'
import vuetify from 'plugins/vuetify'

import 'css/app.scss'

import 'regenerator-runtime/runtime'
import { i18n } from './trans'
import AppMessenger from './plugins/messenger'
import App from './App'

import { NewStore, storeOptions } from './store'
const store = NewStore(storeOptions)

import NewRouter from './router'
const router = NewRouter(store)

import filters from 'plugins/filter'
Vue.use(filters)

import websocket from 'plugins/websocket'
Vue.use(websocket, {
  namespace: 'ws',
  store: store
})

import snackbar from 'plugins/snackbar'
Vue.use(snackbar, {
  store: store
})

import loader from 'plugins/loader'
Vue.use(loader, {
  store: store
})

import { MESSAGES_READ, CONVERSATION_SEND_MESSAGE } from 'actions/conversations'

Vue.use(AppMessenger, {
  namespace: 'conversations',
  mutations: {
    MESSAGES_READ
  },
  actions: {
    CONVERSATION_SEND_MESSAGE
  },
  store
})

new Vue({
  el: '#app',
  router,
  i18n,
  store,
  vuetify,
  websocket,
  AppMessenger,
  render: (h) => h(App)
})
