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

import userRepository from 'repos/user'
import axios from 'repos/repository'

const state = {
  email: localStorage.getItem('user-email') || '',
  token: localStorage.getItem('token') || '',
  status: ''
}

const getters = {
  isAuthenticated: (state) => !!state.email && !!state.token,
  authStatus: (state) => state.status
}

const actions = {
  [AUTH_REQUEST]: ({ commit }, user) => {
    commit(AUTH_REQUEST)
    return userRepository.login(user).then((response) => {
      commit(AUTH_SUCCESS, response.data)
    })
  },
  [AUTH_SIGNUP]: ({ commit }, user) => {
    commit(AUTH_SIGNUP)
    return userRepository.create(user).then(() => {
      commit(AUTH_SIGNUP_SUCCESS)
    })
  },
  [AUTH_CHANGE_PASSWORD]: ({ commit }, user) => {
    commit(AUTH_CHANGE_PASSWORD)
    return userRepository.changePassword(user).then(() => {
      commit(AUTH_CHANGE_PASSWORD_SUCCESS)
    })
  },
  [AUTH_LOGOUT]: async ({ commit }) => {
    commit(AUTH_LOGOUT)
    return userRepository
      .logout()
      .then(() => {
        commit(AUTH_CHANGE_PASSWORD_SUCCESS)
      })
      .catch(console.log)
  },
  [AUTH_ERROR]: ({ commit }) => {
    commit(AUTH_LOGOUT)
    commit(AUTH_ERROR)
  }
}

const mutations = {
  [AUTH_REQUEST]: (state) => {
    state.status = 'loading'
  },
  [AUTH_SUCCESS]: (state, data) => {
    axios.defaults.headers.common.Authorization = data.token
    state.email = data.email
    state.token = data.token
    localStorage.setItem('user-email', data.email)
    localStorage.setItem('token', data.token)
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
  [AUTH_ERROR]: (state) => {
    state.status = 'auth-error'
  },
  [AUTH_LOGOUT]: (state) => {
    localStorage.removeItem('user-email')
    localStorage.removeItem('token')
    delete axios.defaults.headers.common.Authorization
    state.email = null
    state.token = null
    state.status = 'logged-out'
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
