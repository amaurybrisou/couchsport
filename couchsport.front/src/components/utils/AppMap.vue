<template>
  <div v-show="show">
    <div v-show="help" slot="help">
      <v-subheader color="info">
        <v-icon color="info"> mdi-help </v-icon>
        {{ help }}
      </v-subheader>
    </div>
    <l-map
      ref="map"
      :zoom="zoom"
      :center="center"
      :max-bounds="maxBounds"
      :now-wrap="true"
      :style="`height:${height};width:${width};`"
    >
      <span v-if="value.length > 0">
        <l-marker v-for="(m, idx) in value" :key="idx" :lat-lng="m" />
      </span>
      <l-tile-layer :url="url" :attribution="attribution" />
    </l-map>
    <slot v-if="!valid" name="errors">
      <v-alert :color="errorColor">
        {{ error }}
      </v-alert>
    </slot>
  </div>
</template>

<script>
  import { LMap, LMarker, LTileLayer } from 'vue2-leaflet'

  export default {
    name: 'AppMap',
    components: { LMap, LMarker, LTileLayer },
    props: {
      help: { type: String, default: '', required: false },
      zoom: {
        type: Number,
        default: 5,
        required: false,
        validator: (v) => !!v && v >= 0
      },
      height: {
        type: String,
        default: '50vh',
        required: false,
        validator: (v) => !!v && v !== ''
      },
      width: {
        type: String,
        default: '30vw',
        required: false,
        validator: (v) => !!v && v !== ''
      },
      center: {
        value: [Array, Object],
        default: () => [0, 0],
        required: false,
        validator: (m) => {
          if (!(m.lat instanceof Number) && m.lat < -90 && m.lat > 90) {
            this.error = this.errors.invalid
          }
          if (!(m.lng instanceof Number) && m.lng < -180 && m.lng > 180) {
            this.error = this.errors.invalid
          }
          return true
        }
      },
      maxBounds: {
        type: Array,
        default: () => [
          [-90, -180],
          [90, 180]
        ],
        required: false
      },
      errorColor: {
        type: String,
        default: 'warning'
      },
      url: {
        type: String,
        default: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        required: false
      },
      attribution: {
        type: String,
        default:
          '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        required: false
      },
      value: {
        type: Array,
        default: () => [],
        required: true
      },
      max: {
        type: Number,
        default: null,
        required: false
      },
      min: {
        type: Number,
        default: 0,
        required: false
      },
      errors: {
        type: Object,
        default() {
          return {
            too_much: `too much markers : ${this.max}`,
            too_few: `too much markers : ${this.min}`,
            invalid: `invalid marker`
          }
        },
        required: false
      },
      show: { type: Boolean, default: true }
    },
    data() {
      return {
        error: this.value.length ? '' : this.errors.too_few,
        marker: {}
      }
    },
    computed: {
      valid() {
        return !this.error
      }
    },
    watch: {
      show(v) {
        var that = this
        setTimeout(function () {
          !!v && that.map.invalidateSize()
        }, 200)
      }
    },
    mounted() {
      var that = this
      this.$nextTick(function () {
        that.map = this.$refs.map.mapObject
        that.map.on('click', that.validate)
        that.map.on('contextmenu', that.validate)
        this.map.invalidateSize()
      })
    },
    methods: {
      reset() {
        this.value = []
      },
      validate({ latlng, type }) {
        switch (type) {
          case 'click':
            this.value.push(latlng)
            break
          case 'contextmenu':
            this.value.pop()
            break
        }

        this.error = null
        if (isNaN(latlng.lat) || latlng.lat < -90 || latlng.lat > 90) {
          this.error = this.errors.invalid
          return false
        }
        if (isNaN(latlng.lng) || latlng.lng < -180 || latlng.lng > 180) {
          this.error = this.errors.invalid
          return false
        }

        this.error =
          this.value.length < this.min ? this.errors.too_few : this.error
        this.error =
          this.value.length > this.max ? this.errors.too_much : this.error

        this.$emit('input', this.value)
        return true
      }
    }
  }
</script>
