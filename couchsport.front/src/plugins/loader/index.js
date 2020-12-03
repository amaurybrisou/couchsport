import LoaderComponent from './AppLoader'
import { TOGGLE_LOADER } from './action'
import loaderModule from './module'

const Loader = {
  install(Vue, options) {
    const isDef = (v) => v !== undefined
    if (!isDef(options.store)) throw new Error('a vuex store is required')

    const c = options.component || LoaderComponent
    const namespace = options.namespace || 'loader'
    const e = options.event || namespace + ':toggle'

    const temporaryMixin = {
      created: function () {
        options.store.subscribe((mutation, state) => {
          if (mutation.type.includes(namespace + '/')) {
            this.$emit(e, state[namespace])
          }
        })
      }
    }

    c.mixins = [...(c.mixins || []), temporaryMixin]

    Vue.component(c.name, c)

    options.store.registerModule(namespace, loaderModule)

    Vue.prototype['$' + namespace] = function (payload) {
      options.store.dispatch(namespace + '/' + TOGGLE_LOADER, payload)
    }
  }
}

if (typeof window !== 'undefined' && window.Vue) {
  Vue.use(AppMessenger)
}

export default Loader
