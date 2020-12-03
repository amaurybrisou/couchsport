import {
  GET_PAGES,
  GOT_PAGES,
  GET_PAGE,
  GOT_PAGE,
  SAVE_PAGE,
  PAGE_SAVED,
  PAGE_ADD_IMAGE,
  MODIFY_IMAGE_ALT,
  PAGE_DELETE_IMAGE,
  PAGE_IMAGE_DELETED,
  NEW_PAGE,
  NEW_PAGE_SAVED,
  PUBLISH_PAGE,
  PAGE_PUBLISHED,
  DELETE_PAGE,
  PAGE_DELETED,
  EDIT_PAGE,
  CANCEL_EDIT_PAGE,
  MODIFY_PAGE,
  REMOVE_ACTIVITY
} from '../actions/pages'

import pagesRepo from '../../repositories/pages'
import imagesRepo from '../../repositories/images'

import Vue from 'vue'

const state = {
  status: '',
  pages: [],
  edited_page: { activities: [], images: [] }
}

const getters = {}

const actions = {
  [GET_PAGES]: ({ commit }) => {
    commit(GET_PAGES)
    return pagesRepo.mines().then(({ data }) => {
      commit(GOT_PAGES, data)
    })
  },
  [GET_PAGE]: ({ commit }, params) => {
    commit(GET_PAGE)
    return pagesRepo.get(params).then((response) => {
      commit(GOT_PAGE, response.data)
      return response
    })
  },
  [SAVE_PAGE]: ({ commit }, what) => {
    if (['edit', 'new'].indexOf(what) < 0) throw new Error('unknow method')

    what === 'edit' && commit(SAVE_PAGE)
    what === 'new' && commit(NEW_PAGE)

    return pagesRepo[what](state.edited_page).then(({ data }) => {
      what === 'edit' && commit(PAGE_SAVED, data)
      what === 'new' && commit(NEW_PAGE_SAVED, data)
      commit(CANCEL_EDIT_PAGE)
    })
  },
  [DELETE_PAGE]: ({ commit }, page) => {
    commit(DELETE_PAGE)
    return pagesRepo
      .delete(page)
      .then((onfulfilled) => onfulfilled && commit(PAGE_DELETED, page.id))
  },
  [PAGE_DELETE_IMAGE]: ({ commit }, imageidX) => {
    commit(PAGE_DELETE_IMAGE)
    if (!state.edited_page.images[imageidX].id) {
      return commit(PAGE_IMAGE_DELETED, imageidX)
    }
    return imagesRepo
      .delete(state.edited_page.images[imageidX])
      .then(commit(PAGE_IMAGE_DELETED, imageidX))
  },
  [PUBLISH_PAGE]: ({ commit }, page) => {
    commit(PUBLISH_PAGE, page.public)
    return pagesRepo.publish(page).then(commit(PAGE_PUBLISHED, page))
  }
}

const mutations = {
  [GET_PAGES]: (state) => {
    state.status = 'getting_pages'
  },
  [GOT_PAGES]: (state, pages) => {
    state.status = 'got_pages'
    state.pages = pages
  },

  [GET_PAGE]: (state) => {
    state.status = 'getting_page'
  },
  [GOT_PAGE]: (state, page) => {
    if (!Array.isArray(page)) {
      page = [page]
    }
    for (let pageIndex in page) {
      let i = state.pages.findIndex((p) => page[pageIndex].id === p.id)
      if (i < 0) state.pages.push(page[pageIndex])
      state.pages.splice(i, 1, page[pageIndex])
    }
    state.status = 'got_page'
  },

  [SAVE_PAGE]: (state) => {
    state.status = 'saving_page'
  },
  [PAGE_SAVED]: (state, { Images, id }) => {
    for (var i = 0; i < state.pages.length; i++) {
      let p = state.pages[i]
      if (p.id === id) {
        Vue.set(state.pages[i], 'Images', Images)
        break
      }
    }
    state.status = 'page_saved'
  },

  [NEW_PAGE]: (state) => {
    state.status = 'saving_new_page'
  },
  [NEW_PAGE_SAVED]: (state, page) => {
    state.pages.push(page)
    state.status = 'new_page_saved'
  },

  [EDIT_PAGE]: (state, page_id) => {
    state.status = 'editing_page'
    for (var i = 0; i < state.pages.length; i++) {
      let p = state.pages[i]
      if (p.id === page_id) {
        Vue.set(state, 'edited_page', state.pages[i])
        break
      }
    }
  },
  [CANCEL_EDIT_PAGE]: (state) => {
    Vue.set(state, 'edited_page', {
      activities: [],
      images: [],
      couch_number: 0
    })
    state.status = 'edit_page_canceled'
  },

  [MODIFY_PAGE]: (state, payload) => {
    state.status = 'modifying_page'
    state.edited_page = { ...state.edited_page, ...payload }
    state.status = 'page_modified'
  },
  [REMOVE_ACTIVITY]: (state, activity) => {
    state.status = 'removing_activity'
    state.edited_page.activities = state.edited_page.activities.filter(
      (a) => activity.id !== a.id
    )
    state.status = 'activity_removed'
  },
  [PAGE_ADD_IMAGE]: (state, image) => {
    state.edited_page.images.push(image)
    // state.status = "page_adding_photo";
  },
  [PAGE_DELETE_IMAGE]: (state) => {
    state.status = 'page_deleting_image'
  },
  [MODIFY_IMAGE_ALT]: (state, { idx, value }) => {
    Vue.set(state.edited_page.images[idx], 'alt', value)
  },
  [PAGE_IMAGE_DELETED]: (state, imageidX) => {
    state.edited_page.images = state.edited_page.images.filter((i, j) =>
      j !== imageidX ? i : null
    )
    state.status = 'page_photo_deleteed'
  },

  [DELETE_PAGE]: (state) => {
    state.status = 'removing_page'
  },
  [PAGE_DELETED]: (state, pageid) => {
    state.pages = state.pages.filter((p) => pageid !== p.id)
    state.status = 'page_removed'
  },

  [PUBLISH_PAGE]: (state, _public) => {
    state.status = (_public ? '' : 'un') + 'publishing_page'
  },
  [PAGE_PUBLISHED]: (state, { id, public: _public }) => {
    for (var i in state.pages) {
      let p = state.pages[i]
      if (p.id === id) {
        Vue.set(state.pages[i], 'public', _public)
        break
      }
    }
    state.status = 'page_removed'
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
