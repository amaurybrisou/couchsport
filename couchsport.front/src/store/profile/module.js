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

import pages from 'store/pages'
import conversations from 'store/conversations'

import profileRepo from 'repos/profile'
import activityRepo from 'repos/activity'
import languageRepo from 'repos/language'

import axios from 'repos/repository'

import { AUTH_LOGOUT } from 'store/auth/actions'

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
    profileRepo.get().then((response) => {
      commit(PROFILE_SUCCESS, response.data)
    })
  },
  [SAVE_PROFILE]: ({ commit }) => {
    commit(SAVE_PROFILE)
    return profileRepo.update(state.profile).then((response) => {
      commit(PROFILE_SUCCESS, response.data)
    })
  },
  [SET_ACTIVITIES]: ({ commit }) => {
    if (sessionStorage.activities) {
      const activities = JSON.parse(sessionStorage.activities)
      commit(SET_ACTIVITIES, activities)
      return
    }
    commit(GET_ACTIVITIES)
    activityRepo.all().then((response) => {
      commit(SET_ACTIVITIES, response.data)
    })
  },
  [SET_LANGUAGES]: ({ commit }) => {
    if (sessionStorage.languages) {
      const languages = JSON.parse(sessionStorage.languages)
      commit(SET_LANGUAGES, languages)
      return
    }
    commit(GET_LANGUAGES)
    languageRepo.all().then((response) => {
      commit(SET_LANGUAGES, response.data)
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
