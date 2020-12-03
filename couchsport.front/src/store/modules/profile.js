import {
  MODIFY_PROFILE,
  SAVE_PROFILE,
  PROFILE_ERROR,
  PROFILE_REQUEST,
  PROFILE_SUCCESS,
  SET_ACTIVITIES,
  GET_ACTIVITIES,
  GET_LANGUAGES,
  SET_LANGUAGES
} from 'actions/profile'

import pages from './pages'
import conversations from './conversations'

import profileRepo from 'repos/profile'
import activityRepo from 'repos/activity.js'
import languageRepo from 'repos/language.js'

import axios from 'repos/repository.js'

import { AUTH_LOGOUT } from 'actions/auth'

const state = {
  status: '',
  profile: {
    locale: 'fr'
  },
  activities: [],
  languages: []
}

const getters = {
  getProfile: (state) => state.profile,
  isProfileLoaded: (state) => !!state.profile.id,
  getLocale: (state) => state.profile.locale
}

const actions = {
  [PROFILE_REQUEST]: ({ commit }) => {
    commit(PROFILE_REQUEST)
    profileRepo.get().then(({ data }) => {
      commit(PROFILE_SUCCESS, data)
    })
  },
  [SAVE_PROFILE]: ({ commit }) => {
    commit(SAVE_PROFILE)
    return profileRepo.update(state.profile).then(({ data }) => {
      commit(PROFILE_SUCCESS, data)
    })
  },
  [GET_ACTIVITIES]: ({ commit, dispatch }) => {
    if (sessionStorage.activities) {
      const activities = JSON.parse(sessionStorage.activities)
      return commit(SET_ACTIVITIES, activities)
    }
    return dispatch(SET_ACTIVITIES)
  },
  [SET_ACTIVITIES]: ({ commit }) => {
    return activityRepo.all().then(({ data }) => {
      commit(SET_ACTIVITIES, data)
      return data
    })
  },
  [GET_LANGUAGES]: ({ commit, dispatch }) => {
    if (sessionStorage.languages) {
      const languages = JSON.parse(sessionStorage.languages)
      return commit(SET_LANGUAGES, languages)
    }
    return dispatch(SET_LANGUAGES)
  },
  [SET_LANGUAGES]: ({ commit }) => {
    return languageRepo.all().then(({ data }) => {
      commit(SET_LANGUAGES, data)
      return data
    })
  }
}

const mutations = {
  [GET_LANGUAGES]: (state) => {
    state.status = 'loading_languages'
  },
  [SET_LANGUAGES]: (state, languages) => {
    state.languages = languages
    sessionStorage.removeItem('languages')
    sessionStorage.setItem('languages', JSON.stringify(state.languages))
    state.status = 'languages_loaded'
  },
  [GET_ACTIVITIES]: (state) => {
    state.status = 'loading_activities'
  },
  [SET_ACTIVITIES]: (state, activities) => {
    state.activities = activities
    sessionStorage.removeItem('activities')
    sessionStorage.setItem('activities', JSON.stringify(state.activities))
    state.status = 'activities_loaded'
  },
  [PROFILE_REQUEST]: (state) => {
    state.status = 'loading'
  },
  [PROFILE_SUCCESS]: (state, profile) => {
    state.profile = { ...state.profile, ...profile }
    state.status = 'success'
  },
  [PROFILE_ERROR]: (state) => {
    state.status = 'error'
  },
  [MODIFY_PROFILE]: (state, payload) => {
    if (Object.hasOwnProperty('locale')) {
      axios.defaults.headers.common['Accept-Language'] = payload.locale
    }
    state.profile = { ...state.profile, ...payload }
    state.status = 'modified_success'
  },
  [SAVE_PROFILE]: (state) => {
    state.status = 'loading'
  },
  [AUTH_LOGOUT]: (state) => {
    state.profile = {}
  }
}

export default {
  state,
  getters,
  actions,
  mutations,
  modules: {
    pages,
    conversations
  }
}
