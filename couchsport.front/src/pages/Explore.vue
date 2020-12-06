<template>
  <v-container fluid fill-height pa-0>
    <v-layout column>
      <v-flex>
        <v-toolbar flat color="primary">
          <v-toolbar-title class="hidden-xs-only">
            {{ $t('p.explore.sb_title') | capitalize }}
          </v-toolbar-title>
          <v-autocomplete
            v-model="select"
            :data-intro="$t('p.explore.help.first_step')"
            data-step="1"
            :loading="loading"
            :items="items"
            :search-input.sync="search"
            return-object
            item-text="name"
            cache-items
            class="mx-3"
            flat
            dense
            hide-no-data
            hide-details
            background-color="secondary"
            :label="$t('p.explore.sb_placeholder') | capitalize"
            solo
            :menu-props="{ zIndex: '2000' }"
            @input="filterMarkers('spots')"
          />
          <v-btn
            icon
            @click="
              select = null
              filterMarkers('spots')
            "
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-toolbar>
      </v-flex>
      <v-flex xs12>
        <l-map
          ref="map"
          :data-intro="$t('p.explore.help.second_step')"
          data-step="2"
          data-position="top"
          :zoom="mapConfig.zoom"
          :center="mapConfig.center"
          :max-bounds="mapConfig.maxBounds"
          :no-wrap="mapConfig.noWrap"
          @ready="map = $refs.map.mapObject"
        >
          <l-tile-layer
            :url="mapConfig.url"
            :attribution="mapConfig.attribution"
            :no-wrap="mapConfig.noWrap"
          >
          </l-tile-layer>
          <l-marker
            v-for="item in layers[layer].markers"
            :key="item.id"
            :lat-lng="item.latlng"
            :name="item.name"
            :draggable="false"
            :visible="item.show"
          >
            <l-popup
              v-if="item.popup.id"
              :lat-lng="item.popup.latlng"
              :options="layers[layer].popupOptions"
            >
              <marker-popup
                :id="item.popup.id"
                :name="item.popup.name"
                :desc="item.popup.desc"
                :image="item.popup.image"
                :url="item.popup.url"
                :activities="item.popup.activities"
              />
            </l-popup>
          </l-marker>
        </l-map>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import { LMap, LTileLayer, LPopup, LMarker } from 'vue2-leaflet'

  import MarkerPopup from 'components/explore/MarkerPopup'
  import { mapActions, mapState } from 'vuex'
  import { GET_PAGE } from 'store/pages/actions'
  // import IntroJS from 'mixins/intro'

  export default {
    name: 'Explore',
    components: { LMap, LMarker, LPopup, LTileLayer, MarkerPopup },
    // mixins: [IntroJS],
    data() {
      return {
        loading: false,
        items: [],
        search: null,
        select: null,
        map: null,
        autocompleteFeed: [],
        mapConfig: {
          zoom: 5,
          center: {
            lat: 47.41322,
            lng: -1.219482
          },
          // maxBounds: [[-120, -210], [120, 210]],
          noWrap: true,
          url: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
          attribution:
            '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        },
        layer: 'spots',
        layers: {
          spots: {
            id: 0,
            name: 'Spots',
            active: true,
            markers: [],
            popupImage: {
              height: '255px',
              width: '255px'
            },
            popupOptions: {
              maxHeight: 450,
              maxWidth: 450
            }
          }
        }
      }
    },
    computed: {
      ...mapState({ pages: (state) => state.profile.pages.pages })
    },
    watch: {
      search(val) {
        val && val !== this.select && this.querySelections(val)
      },
      pages: {
        handler() {
          this.extractAutoCompleteItems()
          this.initLayers()
        }
      }
    },
    created() {
      this.GET_PAGE()
    },
    mounted() {
      if (this.$route && this.$route.query.zoom)
        this.mapConfig.zoom = this.$route.query.zoom
      if (this.$route && this.$route.query.lat && this.$route.query.lng) {
        this.mapConfig.center.lat = this.$route.query.lat
        this.mapConfig.center.lng = this.$route.query.lng
      } else {
        if (navigator.geolocation) {
          var self = this
          navigator.geolocation.getCurrentPosition(function (position) {
            self.mapConfig.center = {
              lat: position.coords.latitude,
              lng: position.coords.longitude
            }
          })
        }
      }

      // var that = this
      // this.help.setOption('doneLabel', this.$t('help.next_page'))
      // this.help.oncomplete(function () {
      //   that.$router.push({
      //     name: 'page-details',
      //     params: { page_name: 'random' }
      //   })
      // })
      // this.help.start()
    },

    methods: {
      ...mapActions('pages', [GET_PAGE]),
      async initLayers() {
        this.initLayer(this.layer)
      },
      initLayer(layer) {
        for (let index in this.pages) {
          const p = this.pages[index]
          if (!p.public) continue

          this.layers[layer].markers.splice(index, 1, {
            id: p.id,
            name: p.name,
            tags:
              p.activities instanceof Array
                ? [p.name].concat(p.activities.map((e) => e.name))
                : [],
            type: 'marker',
            latlng: [p.lat, p.lng],
            show: true,
            popup: {
              id: p.id,
              url: '/' + this.$i18n.locale + '/pages/' + p.name,
              name: p.name,
              image:
                p.images.length > 0
                  ? {
                      url: p.images[0].url,
                      alt: p.images[0].alt,
                      width: this.layers[layer].popupImage.width,
                      height: this.layers[layer].popupImage.height
                    }
                  : {},
              activities: p.activities == null ? {} : p.activities,
              desc: p.description
            }
          })
        }
      },
      filterMarkers(layer) {
        for (let index in this.layers[layer].markers) {
          const m = this.layers[layer].markers[index]

          if (this.select == null) {
            m.show = true
            this.map.setZoom(2)
            continue
          }
          if (m.tags.length === 0) {
            m.show = false
            continue
          }

          for (let j in m.tags) {
            const markerTag = m.tags[j]
            if (
              (markerTag || '')
                .toLowerCase()
                .indexOf((this.select.name || '').toLowerCase()) > -1
            ) {
              m.show = true
              this.map.setZoom(2)
              break
            }
            m.show = false
          }
        }
      },
      extractAutoCompleteItems() {
        if (this.pages.length > 0) {
          var that = this
          this.pages.forEach((page) => {
            if (page.public) {
              that.autocompleteFeed = that.autocompleteFeed
                .concat(page.activities || [])
                .concat(page)
            }
          })
        }
      },
      querySelections(v) {
        this.items = this.autocompleteFeed.filter((e) => {
          return (e.name || '').toLowerCase().includes(v)
        })
      }
    }
  }
</script>
