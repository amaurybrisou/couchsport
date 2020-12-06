import {
  AUTH_REQUEST,
  AUTH_SIGNUP,
  AUTH_SIGNUP_SUCCESS,
  AUTH_ERROR,
  AUTH_SUCCESS,
  AUTH_LOGOUT,
  AUTH_CHANGE_PASSWORD,
  AUTH_CHANGE_PASSWORD_SUCCESS
} from 'store/auth/actions'

import { PROFILE_REQUEST } from 'store/profile/actions'

import userRepository from 'repos/user'
import axios from 'repos/repository'

const state = {
  email: localStorage.getItem('user-email') || '',
  token: localStorage.getItem('token') || '',
  status: '',
  errors: []
}

const getters = {
  isAuthenticated: (state) => !!state.email && !!state.token,
  authStatus: (state) => state.status,
  errors: (state) => state.errors
}

const actions = {
  [AUTH_REQUEST]: async ({ commit, dispatch }, user) => {
    commit(AUTH_REQUEST)
    const response = await userRepository.login(user)

    if (response.status === 200) {
      commit(AUTH_SUCCESS, response)
      dispatch(PROFILE_REQUEST)
      return response.data
    }
  },
  [AUTH_SIGNUP]: async ({ commit }, user) => {
    commit(AUTH_SIGNUP)
    const response = await userRepository.create(user)
    if (response.status === 200) {
      commit(AUTH_SIGNUP_SUCCESS, response)
      return response.data
    }
  },
  [AUTH_CHANGE_PASSWORD]: async ({ commit }, user) => {
    commit(AUTH_CHANGE_PASSWORD)
    const response = await userRepository.changePassword(user)

    if (response.status === 200) {
      commit(AUTH_CHANGE_PASSWORD_SUCCESS)
      return response.data
    }
  },
  [AUTH_LOGOUT]: async ({ commit }) => {
    commit(AUTH_LOGOUT)
    return userRepository.logout()
  },
  [AUTH_ERROR]: ({ commit }) => {
    commit(AUTH_LOGOUT)
  }
}

const mutations = {
  [AUTH_REQUEST]: (state) => {
    state.status = 'loading'
  },
  [AUTH_SUCCESS]: (state, resp) => {
    axios.defaults.headers.common.Authorization = resp.data.token
    localStorage.setItem('user-email', resp.data.email)
    localStorage.setItem('token', resp.data.token)
    state.email = resp.data.email
    state.status = 'success'
  },
  [AUTH_CHANGE_PASSWORD]: (state) => {
    state.status = 'changing-password'
  },
  [AUTH_CHANGE_PASSWORD_SUCCESS]: (state) => {
    state.status = 'password-changed'
  },
  [AUTH_SIGNUP]: (state) => {
    state.status = 'signing-up'
  },
  [AUTH_SIGNUP_SUCCESS]: (state) => {
    state.status = 'signup-success'
  },
  [AUTH_ERROR]: (state, error) => {
    state.status = 'error'
    state.errors = [error]
    state.email = null
    localStorage.removeItem('user-email')
  },
  [AUTH_LOGOUT]: (state) => {
    localStorage.removeItem('user-email')
    delete axios.defaults.headers.common.Authorization
    state.email = null
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
