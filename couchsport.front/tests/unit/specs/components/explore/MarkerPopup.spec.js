import { shallowMount, createLocalVue } from '@vue/test-utils'
import Vuetify from 'vuetify'
import { i18n } from '@/trans'

import MarkerPopup from 'components/explore/MarkerPopup.vue'

const localVue = createLocalVue()
localVue.use(i18n)

describe('MarkerPopup.vue', () => {
  let deps = {}
  let wrapper
  beforeEach(() => {
    deps = {
      localVue,
      i18n,
      vuetify: new Vuetify()
    }
  })

  it('should render correctly', () => {
    deps.propsData = {
      name: 'TEST NAME',
      desc: 'TEST DESC',
      image: {},
      id: 1,
      url: 'https://picsum.photos/200/300',
      activities: ['Surf']
    }

    wrapper = shallowMount(MarkerPopup, deps)

    expect(wrapper.text()).toContain('TEST NAME')
  })

  it('should render correctly', () => {
    deps.propsData = {
      id: 1
    }

    wrapper = shallowMount(MarkerPopup, deps)

    let spy = jest.spyOn(wrapper.vm.$options.props.id, 'validator')

    try {
      spy()
    } catch (e) {
      expect(e).toEqual(new Error('Id cannot be undefined or below 0'))
    }

    spy.mockRestore()

    spy = jest.spyOn(wrapper.vm.$options.props.id, 'default')

    try {
      spy()
    } catch (e) {
      expect(e).toEqual(new Error('Id cannot be undefined or below 0'))
    }

    spy.mockRestore()
  })
})
