import { createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import { storeOptions, NewStore } from '@/store'

import { AUTH_SIGNUP } from 'store/auth/actions'

import repo from 'repos/repository'

import { cloneDeep } from 'lodash'

const localeVue = createLocalVue()
localeVue.use(Vuex)

jest.mock('repos/repository')

describe('store.auth', () => {
  let store
  beforeEach(function () {
    store = NewStore(cloneDeep(storeOptions))
  })

  it('signups correctly', async () => {
    const user = {
      email: 'amaury.brisou@puzzle.org',
      password: 'TestPassword123'
    }

    repo.put.mockImplementationOnce(() => {
      return Promise.resolve({ status: 200, data: user })
    })

    const response = await store.dispatch(AUTH_SIGNUP, user)
    expect(response).toEqual(user)
    expect(store.state.auth.status).toEqual('signup-success')
  })
})
