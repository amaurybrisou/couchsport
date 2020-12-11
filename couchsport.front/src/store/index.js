import Vue from 'vue'
import Vuex from 'vuex'

import profile from 'store/profile/module'
import auth from 'store/auth/module'

import axios from 'repos/repository'
import NewRouter from '@/router'

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
  const router = NewRouter(store)

  axios.interceptors.response.use(
    function (r) {
      return r
    },
    function (error) {
      if (
        error.response &&
        (error.response.status === 403 || error.response.status === 401)
      ) {
        store
          .dispatch(AUTH_ERROR)
          .catch(() => console.log)
          .finally(
            router.push({
              name: 'login',
              params: {
                locale: store.getters.getLocale
              }
            })
          )
      }

      return Promise.reject(`api_errors.${error.response.data}`)
    }
  )

  return store
}
