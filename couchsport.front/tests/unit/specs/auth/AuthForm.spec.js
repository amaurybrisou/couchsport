import { mount, createLocalVue } from '@vue/test-utils'
import Vuex from 'vuex'
import filters from 'plugins/filter'
import Vuetify from 'vuetify'
import { i18n } from '@/trans'

import { storeOptions, NewStore } from '@/store'
import AuthForm from '@/components/auth/AuthForm.vue'

import { cloneDeep, isEqual } from 'lodash'

const localVue = createLocalVue()
localVue.use(Vuex)
localVue.use(filters)
localVue.use(i18n)

describe('AuthForm.vue', () => {
  let deps = {}
  let wrapper
  let submit

  beforeEach(() => {
    deps = {
      localVue,
      i18n,
      vuetify: new Vuetify(),
      store: new NewStore(cloneDeep(storeOptions))
    }
  })

  afterEach(function () {
    wrapper.destroy()
  })

  it('dispatches and executes "submit" with a valid email/password', async () => {
    const user = {
      email: 'amaury.brisou@puzzledge.org',
      password: 'Abcdefgh1'
    }

    deps.data = () => {
      return {
        user
      }
    }

    deps.propsData = {
      submit: function (pUser) {
        expect(isEqual(pUser, user)).toBeTruthy()
      }
    }

    wrapper = mount(AuthForm, deps)

    submit = jest.spyOn(wrapper.vm, 'submitForm')
    await localVue.nextTick()

    await wrapper.find('[name=submit]').trigger('click')
    expect(wrapper.vm.$data.valid).toBe(true)
    expect(submit).toHaveBeenCalledTimes(1)
    expect(wrapper.emitted().submit).toBeTruthy()
  })

  it('dispatches "submit" with a invalid password', async () => {
    deps.data = () => {
      return {
        user: {
          email: 'amaury.brisou@puzzledge.org',
          password: 'Abcdefgh'
        }
      }
    }
    wrapper = mount(AuthForm, deps)

    submit = jest.spyOn(wrapper.vm, 'submitForm')

    await localVue.nextTick()
    await wrapper.find('[name=submit]').trigger('click')
    expect(wrapper.vm.$data.valid).toBe(false)
    expect(submit).not.toHaveBeenCalledTimes(1)
  })

  it('dispatches "submit" with a invalid email', async () => {
    deps.data = () => {
      return {
        user: {
          email: 'amaury.brisou',
          password: 'Abcdefgh1'
        }
      }
    }
    wrapper = mount(AuthForm, deps)

    submit = jest.spyOn(wrapper.vm, 'submitForm')

    await localVue.nextTick()
    await wrapper.find('[name=submit]').trigger('click')
    expect(wrapper.vm.$data.valid).toBe(false)
    expect(submit).not.toHaveBeenCalledTimes(1)
  })

  it('emits "hide-change-password-dialog" if state.email', async () => {
    storeOptions.modules.auth.state.email = 'amaury.brisou@gmail.com'
    storeOptions.modules.auth.state.token = 'ae6d17f1-2dc8-4754'
    storeOptions.modules.profile.state.profile.id = '85480'
    deps.store = new NewStore(cloneDeep(storeOptions))

    wrapper = await mount(AuthForm, deps)

    expect(wrapper.vm.isAuthenticated).toBeTruthy()
    expect(wrapper.vm.isProfileLoaded).toBeTruthy()
    expect(wrapper.vm.email).toEqual('amaury.brisou@gmail.com')
    const cancel = wrapper.find('[name=cancel]')
    expect(cancel.isVisible()).toBeTruthy()

    await cancel.trigger('click')

    expect(wrapper.emitted()['hide-change-password-dialog']).toBeTruthy()
  })
})
