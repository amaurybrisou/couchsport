import moment from 'moment'
import Vue from 'vue'

const AppFilters = {
  install: function (Vue) {
    Vue.filter('formatDate', function (value, format) {
      if (value) {
        return moment(String(value)).format(format || 'MM/DD/YYYY hh:mm')
      }
    })

    Vue.filter('shorten', function (value, max = 10) {
      if (value) {
        return value.slice(0, max) + '...'
      }
    })

    Vue.filter('capitalize', function (value) {
      if (value) {
        return value.charAt(0).toUpperCase() + value.slice(1)
      }
    })
  }
}

if (typeof window !== 'undefined' && window.Vue) {
  Vue.use(AppFilters)
}

export default AppFilters.install
