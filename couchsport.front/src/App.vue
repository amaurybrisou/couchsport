<template>
  <v-app>
    <v-carousel
      v-if="
        [
          `/${$i18n.locale}`,
          `/${$i18n.locale}/`,
          `/${$i18n.locale}/about`,
          `/${$i18n.locale}/signup`,
          `/${$i18n.locale}/login`
        ].indexOf($route.path) >= 0
      "
      class="background-carousel"
      height="100vh"
      hide-delimiters
      cycle
      :show-arrows="false"
    >
      <v-carousel-item
        v-for="(bg, i) in backgrounds"
        :key="i"
        class="carousel-item"
        :src="bg.src"
      >
      </v-carousel-item>
    </v-carousel>
    <app-nav />
    <v-main>
      <transition name="fade" mode="out-in">
        <router-view />
      </transition>
    </v-main>
    <app-snack-bar />
    <app-loader />
  </v-app>
</template>

<script>
  import WebFontLoader from 'webfontloader'
  import AppNav from 'components/nav/AppNav'

  import { PROFILE_REQUEST } from 'store/profile/actions'

  import L from 'leaflet'
  delete L.Icon.Default.prototype._getIconUrl

  L.Icon.Default.mergeOptions({
    iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
    iconUrl: require('static/img/marker-icon.png'),
    shadowUrl: require('leaflet/dist/images/marker-shadow.png')
  })

  export default {
    name: 'App',
    components: {
      'app-nav': AppNav
    },
    data() {
      return {
        backgrounds: [
          { src: require('static/img/bg.jpg') },
          { src: require('static/img/bg1.jpg') },
          { src: require('static/img/bg2.jpg') },
          { src: require('static/img/bg3.jpg') }
        ]
      }
    },
    created: function () {
      if (this.$store.getters.isAuthenticated) {
        this.$store.dispatch(PROFILE_REQUEST)
      }
    },
    mounted() {
      WebFontLoader.load({
        google: {
          families: ['Roboto:100,300,400,500,700,900']
        },
        active: this.setFontLoaded
      })
    },
    methods: {
      setFontLoaded() {
        this.$emit('font-loaded')
      }
    }
  }
</script>

<style lang="scss">
  #app {
    background: none;

    .background-carousel {
      position: absolute;
      z-index: -1;
    }
  }

  .highlight-help {
    border-radius: 5px;
  }

  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.3s;
  }
  .fade-enter,
  .fade-leave-to,
  .fade-leave-active {
    opacity: 0;
  }

  .help {
    $tooltip-color: rgba(#fff, 0.9);
    $tooltip-text-color: rgba(#16191b, 1);
    $tooltip-border-color: rgba(#16191b, 0.8);
    border-radius: 10px;
    background-color: $tooltip-color;
    color: $tooltip-text-color;

    .introjs-arrow.top {
      border-bottom-color: $tooltip-color;
    }

    .introjs-arrow.bottom {
      border-top-color: $tooltip-color;
    }

    .introjs-button {
      color: $tooltip-text-color;
      border-color: $tooltip-border-color;
    }

    .introjs-disabled {
      display: none;
    }

    .introjs-prevbutton {
      margin-right: 5px;
    }
  }
</style>
