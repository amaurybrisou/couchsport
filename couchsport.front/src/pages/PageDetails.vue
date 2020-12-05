<template>
  <v-container>
    <div
      class="app-background"
      :style="{ 'background-image': 'url(' + backgroundImage + ')' }"
    />
    <v-layout
      :class="{
        'column mt-1': $vuetify.breakpoint.xsOnly,
        'row wrap pr-2': $vuetify.breakpoint.smAndUp
      }"
    >
      <v-flex xs12 sm6>
        <!-- <v-card flat class="transparent" fill-height></v-card> -->
        <v-card v-if="page" class="flexcard fill-height">
          <v-card-title class="title font-weight-bold pb-0">
            <div class="font-weight-bold">
              {{ page.name }}
            </div>
            <v-spacer />

            <v-tooltip id="step-1" bottom>
              <template v-slot:activator="{ on }">
                <v-icon color="primary" v-on="on">mdi-web</v-icon>
              </template>
              <span>{{ talkedLanguages }}</span>
            </v-tooltip>
            <v-spacer />
            <v-tooltip v-if="page.couch_number > 0" id="step-2" bottom>
              <template v-slot:activator="{ on }">
                <v-chip color="info" text-color="white" small v-on="on">
                  {{ page.couch_number }} {{ $t('p.pd.avail_couch') }}
                </v-chip>
              </template>
              <span>{{ $t('p.pd.guests') }}</span>
            </v-tooltip>
            <div v-else>
              <v-chip color="primary" text-color="white" small>
                {{ $t('p.pd.no_guests') }}
              </v-chip>
            </div>
          </v-card-title>

          <v-card-text class="grow pa-0">
            <v-list avatar>
              <v-list-item v-if="page.activities" id="step-3">
                <v-chip
                  v-for="(a, i) in page.activities"
                  :key="i"
                  color="primary"
                  text-color="white"
                  small
                  class="mr-2"
                >
                  {{ a.name | capitalize }}
                </v-chip>
              </v-list-item>
              <v-divider />
              <v-list-item class="text-break subheading font-weight-regular">
                {{ page.long_description || page.description }}
              </v-list-item>
            </v-list>
          </v-card-text>

          <v-card-actions id="step-4" class="ma-0 pa-0">
            <v-btn
              depressed
              color="primary"
              block
              :disabled="
                message.from_id == message.to_id || page.couch_number == 0
              "
              @click="showContactDialog = true"
            >
              {{ $t('contact') }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
      <v-flex
        :class="{
          'xs12 mt-2': $vuetify.breakpoint.xsOnly,
          'sm6 pl-2': $vuetify.breakpoint.smAndUp
        }"
      >
        <v-card pa-5>
          <l-map
            ref="map"
            :zoom="mapConfig.zoom"
            :center="mapConfig.center"
            :max-bounds="mapConfig.maxBounds"
            :no-wrap="mapConfig.noWrap"
            style="height: 45vh; width: 100%"
            @ready="map = $refs.map.mapObject"
          >
            <l-tile-layer
              :url="mapConfig.url"
              :attribution="mapConfig.attribution"
            />
          </l-map>
        </v-card>
      </v-flex>
      <v-flex v-if="page && page.images" wrap xs12>
        <v-card
          v-for="(image, idx) in page.images"
          :key="idx"
          align-content-space-between
          class="rounded mt-2"
        >
          <v-img
            max-height="250px"
            :src="image.url"
            :lazy-src="image.url"
            aspect-ratio="1"
            class="grey lighten-2"
            @click="showImageDialog = true"
          >
            <v-layout
              slot="placeholder"
              fill-height
              align-center
              justify-center
              ma-0
            >
              <v-progress-circular indeterminate color="grey lighten-5" />
            </v-layout>
          </v-img>
        </v-card>

        <v-dialog id="image-dialog" v-model="showImageDialog">
          <v-carousel interval="700000000" height="80vh" hide-delimiters>
            <v-icon
              large
              class="right pa-0 ma-1 close-icon"
              icon
              @click="showImageDialog = false"
            >
              mdi-close
            </v-icon>
            <v-carousel-item
              v-for="(image, i) in page.images"
              :key="i"
              :src="image.url"
            />
          </v-carousel>
        </v-dialog>
      </v-flex>
    </v-layout>
    <v-layout row justify-center>
      <v-dialog
        v-if="page && message.from_id != message.to_id"
        id="contact-dialog"
        v-model="showContactDialog"
        width="500"
      >
        <v-card>
          <v-toolbar color="primary">
            <v-card-title class="title font-weight-regular">
              {{ $t('your') }} {{ $t('_message') }} {{ $t('to') }}
              {{ contactName }}
            </v-card-title>
          </v-toolbar>
          <v-form v-model="messageFormValid">
            <v-card-text>
              <v-text-field
                v-if="!email"
                v-model="message.email"
                name="email"
                :label="$t('email')"
                autocomplete="email"
                :rules="emailRules"
              />
              <v-textarea
                v-model="message.text"
                name="message"
                :label="$t('_message')"
                :rules="textRules"
                row="1"
                maxlength="128"
                no-resize
                @keyup.ctrl.enter="sendMessage"
              />
            </v-card-text>
            <v-card-actions>
              <v-spacer />
              <v-btn color="primary" text @click="showContactDialog = false">
                {{ $t('cancel') }}
              </v-btn>
              <v-btn
                color="primary"
                text
                :disabled="!messageFormValid"
                @click="sendMessage"
              >
                {{ $t('send') }}
              </v-btn>
            </v-card-actions>
          </v-form>
        </v-card>
      </v-dialog>
    </v-layout>
  </v-container>
</template>

<script>
  import { LMap, LTileLayer } from 'vue2-leaflet'

  import { GET_PAGE } from 'actions/pages'

  import { mapState, mapActions } from 'vuex'
  import IntroJS from 'mixins/intro'
  export default {
    name: 'PageDetails',
    components: { LMap, LTileLayer },
    mixins: [IntroJS],
    data() {
      return {
        contactName: '',
        backgroundImage: '',
        page: null,
        message: {
          from_id: this.from_id,
          to_id: null,
          email: this.email,
          text: ''
        },

        messageFormValid: false,
        emailRules: [
          (v) => !!v || this.$t('message.required', ['', this.$t('email')]),
          (v) =>
            /.+@.+/.test(v) || this.$t('message.invalid', [this.$t('email')])
        ],

        textRules: [
          (v) => !!v || this.$t('message.required', ['', this.$t('_message')]),
          (v) => (v && v.length >= 20) || this.$t('message.length_above', [20])
        ],

        showImageDialog: false,
        showContactDialog: false,

        map: null,
        mapConfig: {
          zoom: 11,
          center: [46, -1],
          maxBounds: [
            [-90, -180],
            [90, 180]
          ],
          noWrap: true,
          url: 'http://{s}.tile.osm.org/{z}/{x}/{y}.png',
          attribution:
            '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors',
          showMarkers: true,
          hasSpotMarker: false,
          spotMarker: null
        }
      }
    },
    computed: {
      ...mapState({
        email: (state) => state.auth.email,
        from_id: (state) => state.profile.profile.id
      }),
      talkedLanguages() {
        let that = this
        return this.page.owner.languages
          .reduce(
            function (acc, cur, idx, src) {
              acc.push(cur.name)
              acc.push(idx < src.length - 1 ? that.$t('and') : '')
              return acc
            },
            [this.$t('talk')]
          )
          .join(' ')
      },
      imagesUrl() {
        return this.page.images.map((e) => e.url)
      }
    },
    async created() {
      if (!this.$route.params.page_name) {
        return this.$router.push({ name: 'home' })
      }

      this.page = await this.GET_PAGE({
        name: this.$route.params.page_name,
        profile: true
      })
        .then((response) => {
          if (
            !response.data ||
            !Array.isArray(response.data) ||
            response.data.length === 0
          ) {
            throw 'message.not_found'
          }
          var page = response.data[0]
          this.message.to_id = page.owner_id
          this.message.from_id = this.from_id
          this.message.email = this.email
          this.contactName =
            page.owner.username ||
            page.owner.firstname ||
            page.owner.lastname ||
            page.owner.email

          this.backgroundImage =
            page.images && page.images.length > 0 ? page.images[0].url : ''

          this.map.setView([page.lat, page.lng])
          L.marker([page.lat, page.lng]).addTo(this.map)

          setTimeout(this.showHelp, 100)
          return page
        })
        .catch((err) => {
          this.$router.push({
            name: 'home',
            params: { err: { key: err, args: ['page'] } }
          })
        })
    },
    methods: {
      ...mapActions('pages', [GET_PAGE]),
      async sendMessage() {
        if (!this.message.to_id) return

        this.$loader(true)
        const response = await this.$messenger
          .sendMessage(this.message)
          .catch(function () {
            this.$loader(false)
            this.$snackbar(
              this.$t('message.error_sending', [this.$t('_message')])
            )
          })
        this.showContactDialog = false
        this.$loader(false)

        if (response.data) {
          this.$snackbar(
            this.$t('message.success_sending', [this.$t('_message')])
          )
        }
      },
      showHelp() {
        // this.help.addSteps([
        //   {
        //     element: document.querySelector('#step-1'),
        //     intro: this.$t('p.pd.help.first_step')
        //   },
        //   {
        //     element: document.querySelector('#step-2'),
        //     intro: this.$t('p.pd.help.second_step')
        //   },
        //   {
        //     element: document.querySelector('#step-3'),
        //     intro: this.$t('p.pd.help.third_step'),
        //     position: 'top'
        //   },
        //   {
        //     element: document.querySelector('#step-4'),
        //     intro: this.$t('p.pd.help.fourth_step'),
        //     position: 'top'
        //   }
        // ])
        // this.help.start()
      }
    }
  }
</script>

<style lang="scss">
  #image-dialog,
  #contact-dialog,
  .close-icon {
    z-index: 1100;
  }

  .close-icon {
    position: absolute;
    right: 0;
  }

  .rounded {
    border-radius: 5px;
  }

  .round {
    border-radius: 50%;
  }

  .transparent {
    background-color: rgba(#fff, 0.3);
  }

  .app-background {
    position: absolute;
    background-size: cover;
    opacity: 0.4;
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
  }

  .v-tooltip__content {
    z-index: 10000 !important;
  }

  .flexcard {
    display: flex;
    flex-direction: column;

    .text-break {
      word-break: break-all;
      overflow-y: auto;
    }
  }
</style>
