import Vue from 'vue'
import Vuex from 'vuex'

import profile from 'store/profile'
import auth from 'store/auth'

import axios from 'repos/repository'

import { AUTH_ERROR } from 'store/auth/actions'

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
      let ret = error
      if (error && error.response) ret = error.response.data
      store.commit(AUTH_ERROR, ret)
      // return Promise.reject(ret)
      return error
    }
  )

  return store
}
