import SnackBarComponent from './AppSnackBar'
import { SHOW_SNACKBAR } from './action'
import snackbarModule from './module'

/**
 * options must include a store
 * optional options: {
 *  event: <custom event name> default: 'snackbar:show'
 *  component: <a custom component which has to listen on 'snackbar:show' or his specific event>
 * }
 */
const SnackBar = {
  install(Vue, options) {
    const isDef = (v) => v !== undefined
    if (!isDef(options.store)) throw new Error('a vuex store is required')

    const c = options.component || SnackBarComponent
    const namespace = options.namespace || 'snackbar'
    const e = options.event || 'snackbar:show'

    const temporaryMixin = {
      created: function () {
        options.store.subscribe((mutation, state) => {
          if (mutation.type === namespace + '/' + SHOW_SNACKBAR) {
            this.$emit(e, state[namespace])
          }
        })
      }
    }

    c.mixins = [...(c.mixins || []), temporaryMixin]

    Vue.component(c.name, c)

    options.store.registerModule(namespace, snackbarModule)

    Vue.prototype['$' + namespace] = function (payload) {
      options.store.dispatch(namespace + '/' + SHOW_SNACKBAR, payload)
    }
  }
}

if (typeof window !== 'undefined' && window.Vue) {
  Vue.use(AppMessenger)
}

export default SnackBar
