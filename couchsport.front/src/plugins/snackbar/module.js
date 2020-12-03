import { SHOW_SNACKBAR } from './action'

export default {
  namespaced: true,
  state: {
    text: '',
    color: 'warning',
    timeout: 3000
  },
  mutations: {
    [SHOW_SNACKBAR](state, payload) {
      if (typeof payload === 'string') {
        payload = { text: payload }
      }
      state.text = payload.text
      state.color = payload.color || state.color
      state.timeout = payload.timeout || state.timeout
    }
  },
  actions: {
    [SHOW_SNACKBAR]({ commit }, payload) {
      commit(SHOW_SNACKBAR, payload)
    }
  }
}
