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
} from 'store/profile/actions'

import pages from 'store/pages/module'
import conversations from 'store/conversations/module'

import profileRepo from 'repos/profile'
import activityRepo from 'repos/activity'
import languageRepo from 'repos/language'

import axios from 'repos/repository'

import { AUTH_LOGOUT } from 'store/auth/actions'

const state = {
  status: '',
  profile: {
    locale: localStorage.getItem('locale') || 'fr'
  },
  activities: [],
  languages: []
}

const getters = {
  getProfile: (state) => state.profile,
  isProfileLoaded: (state) => !!state.profile.id,
  getLocale: (state) => localStorage.getItem('locale') || state.profile.locale
}

const actions = {
  [PROFILE_REQUEST]: ({ commit }) => {
    commit(PROFILE_REQUEST)
    return profileRepo.get().then((response) => {
      commit(PROFILE_SUCCESS, response.data)
      return response
    })
  },
  [SAVE_PROFILE]: ({ commit }) => {
    commit(SAVE_PROFILE)
    return profileRepo.update(state.profile).then((response) => {
      commit(PROFILE_SUCCESS, response.data)
      return response
    })
  },
  [SET_ACTIVITIES]: ({ commit }) => {
    if (sessionStorage.activities) {
      const activities = JSON.parse(sessionStorage.activities)
      commit(SET_ACTIVITIES, activities)
      return
    }
    commit(GET_ACTIVITIES)
    return activityRepo.all().then((response) => {
      commit(SET_ACTIVITIES, response.data)
      return response
    })
  },
  [SET_LANGUAGES]: ({ commit }) => {
    if (sessionStorage.languages) {
      const languages = JSON.parse(sessionStorage.languages)
      commit(SET_LANGUAGES, languages)
      return
    }
    commit(GET_LANGUAGES)
    return languageRepo.all().then((response) => {
      commit(SET_LANGUAGES, response.data)
      return response
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
    if (payload.hasOwnProperty('locale')) {
      axios.defaults.headers.common['Accept-Language'] = payload.locale
      localStorage.setItem('locale', payload.locale)
    }
    state.profile = { ...state.profile, ...payload }
    state.status = 'success'
  },
  [SAVE_PROFILE]: (state) => {
    state.status = 'loading'
  },
  [AUTH_LOGOUT]: (state) => {
    state.profile = {}
    state.status = 'logout'
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
