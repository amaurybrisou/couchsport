import Vue from 'vue'
import Vuex from 'vuex'
import profile from './modules/profile'
import auth from './modules/auth'
import axios from 'repos/repository'

import { AUTH_ERROR } from 'actions/auth'

Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

export const storeOptions = {
  modules: {
    profile,
    auth
  },
  strict: debug,
  devTools: debug
}

export const NewStore = function (storeOptions) {
  const store = new Vuex.Store(storeOptions)

  axios.interceptors.response.use(
    function (r) {
      return r
    },
    function (error) {
      if (
        error.response &&
        (error.response.status === 403 || error.response.status === 401)
      ) {
        store.dispatch(AUTH_ERROR)
      }
      store.commit(AUTH_ERROR, error)
      let ret = error
      if (error && error.response) ret = error.response.data
      return Promise.reject(ret)
    }
  )

  return store
}
