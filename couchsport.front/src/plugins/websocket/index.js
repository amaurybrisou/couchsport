import VueNativeSock from 'vue-native-websocket'
import WsMutations from './actions'
import WsModule from './modules'

const WebSockets = {
  install: function (Vue, options) {
    const isDef = (v) => v !== undefined
    if (!isDef(options.store)) throw new Error('a vuex store is required')
    const namespace = options.namespace || 'ws'
    const prefix = options.prefix || '/api/ws'
    const protocol = options.ssl ? 'wss' : 'ws'

    options.store.registerModule(namespace, WsModule)

    const url = `${protocol}://${window.location.hostname}:${window.location.port}${prefix}`

    Vue.use(VueNativeSock, url, {
      connectManually: true,
      store: options.store,
      namespace: namespace,
      mutations: WsMutations,
      format: 'json',
      reconnection: true,
      reconnectionAttempts: 5,
      passToStoreHandler: function (eventName, event) {
        if (!eventName.startsWith('SOCKET_')) {
          return
        }
        let method = 'commit'
        let target = eventName.toUpperCase()
        let msg = event
        if (this.format === 'json' && event.data) {
          msg = JSON.parse(event.data)
          if (msg.mutation) {
            target = [msg.namespace || '', msg.mutation]
              .filter((e) => !!e)
              .join('/')
          } else if (msg.action) {
            method = 'dispatch'
            target = [msg.namespace || '', msg.action]
              .filter((e) => !!e)
              .join('/')
          }
        } else {
          target = namespace + '/' + target
        }
        this.store[method](target, msg)
      }
    })

    Vue.directive(namespace, {
      name: 'ws',
      inserted: function (el, binding) {
        if (binding.modifiers.connect) {
          options.store.dispatch(
            `${namespace}/${WsMutations.SOCKET_CONNECT}`,
            `${url}?id=${binding.value}`
          )
        }
      }
    })
  }
}

export default WebSockets
