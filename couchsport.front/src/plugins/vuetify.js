import Vue from 'vue'
import Vuetify from 'vuetify/lib'

import 'vuetify/dist/vuetify.min.css'

Vue.use(Vuetify)

export default new Vuetify({
  theme: {
    themes: {
      light: {
        primary: '#06B998',
        secondary: '#68ccb9',
        accent: '#00bcd4',
        error: '#ff5722',
        warning: '#ff9800',
        info: '#ffc107',
        success: '#607d8b'
      }
    }
  }
})
