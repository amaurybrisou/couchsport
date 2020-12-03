import { mapState } from 'vuex'

// based on article:
// VueJS - Two way data binding and state management with Vuex and strict mode
// https://ypereirareis.github.io/blog/2017/04/25/vuejs-two-way-data-binding-state-management-vuex-strict-mode/

export default function mapStatesTwoWay(namespace, states, updateCb) {
  let mappedStates
  if (namespace) mappedStates = mapState(namespace, states)
  else mappedStates = mapState(states)
  const res = {}
  for (const key in mappedStates) {
    res[key] = {
      set(value) {
        updateCb.call(this, { [key]: value })
      },
      get() {
        return mappedStates[key].call(this)
      }
    }
  }
  return res
}
