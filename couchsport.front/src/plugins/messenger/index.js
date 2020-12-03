import Vue from 'vue'

const AppMessenger = {
  install: function (Vue, options) {
    const isDef = (v) => v !== undefined
    if (!isDef(options.store)) throw new Error('a vuex store is required')
    const namespace = (options.namespace || 'conversations') + '/'

    Vue.prototype.$messenger = {
      mutations: options.mutations,
      actions: options.actions,
      setMessagesRead(conversationIDX) {
        options.store.commit(
          namespace + this.mutations.MESSAGES_READ,
          conversationIDX
        )
      },
      sendMessage(m) {
        return options.store.dispatch(
          namespace + this.actions.CONVERSATION_SEND_MESSAGE,
          m
        )
      }
    }
  }
}

export default AppMessenger

if (typeof window !== 'undefined' && window.Vue) {
  Vue.use(AppMessenger)
}
