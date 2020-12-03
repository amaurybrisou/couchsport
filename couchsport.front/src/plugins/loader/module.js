import { SHOW_LOADER, HIDE_LOADER, TOGGLE_LOADER } from './action'

export default {
  namespaced: true,
  state: {
    message: 'message.stand_by',
    status: false
  },
  mutations: {
    [TOGGLE_LOADER](state, message) {
      if (message) state.message = message
      state.status = !state.status
    },
    [SHOW_LOADER](state, message) {
      if (message) state.message = message
      state.status = true
    },
    [HIDE_LOADER](state, message) {
      if (message) state.message = message
      state.status = false
    }
  },
  actions: {
    [TOGGLE_LOADER]({ commit }, payload) {
      let mutation = TOGGLE_LOADER
      if (typeof payload === 'boolean') {
        mutation = payload ? SHOW_LOADER : HIDE_LOADER
      }
      commit(mutation, payload.message)
    }
  }
}
