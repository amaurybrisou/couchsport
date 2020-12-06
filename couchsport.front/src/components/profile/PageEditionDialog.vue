<template>
  <v-dialog
    v-model="showEditPageDialog"
    fullscreen
    transition="dialog-bottom-transition"
    :overlay="false"
    style="z-index: 600"
    @keydown.esc="cancelEdit()"
  >
    <template v-slot:activator="props">
      <slot name="open-btn" v-bind="props"></slot>
    </template>

    <v-card>
      <v-toolbar dark color="primary">
        <v-toolbar-title>
          <slot name="pageTitle"> {{ $t('edit_page') }} : {{ name }} </slot>
        </v-toolbar-title>
        <v-spacer />
        <v-toolbar-items />
        <v-btn dark text @click.native="validate">
          <slot name="submitText">
            {{ $t('save') }}
          </slot>
        </v-btn>
        <v-btn icon dark @click.native.prevent="cancelEdit()">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-toolbar>
      <template>
        <v-form
          ref="form"
          v-model="rules.valid"
          lazy-validation
          @keypress.enter="validate"
        >
          <v-container
            fluid
            :class="{
              'sm4 px-5 pb-0': $vuetify.breakpoint.smAndUp,
              'xs12 pa-2': $vuetify.breakpoint.xsOnly
            }"
          >
            <v-layout row wrap>
              <v-flex
                :class="{
                  'sm6 pr-1 pb-0': $vuetify.breakpoint.smAndUp,
                  'xs12 pa-2': $vuetify.breakpoint.xsOnly
                }"
              >
                <v-text-field
                  v-model="name"
                  :label="$t('name') | capitalize"
                  :rules="rules['name']"
                  autofocus
                  required
                  @keypress.enter="validate"
                />
                <v-text-field
                  v-model="description"
                  :label="$t('description') | capitalize"
                  :rules="rules['description']"
                  required
                  @keypress.enter="validate"
                />
                <v-textarea
                  v-model="long_description"
                  name="long_description"
                  :rules="rules['long_description']"
                  maxlength="512"
                  :placeholder="$t('p.ped.long_desc_ph') | capitalize"
                  rows="2"
                  no-resize
                  @keypress.ctrl.enter="validate"
                />
                <v-autocomplete
                  v-model="activities"
                  :items="allActivities"
                  :label="$t('activities') | capitalize"
                  :rules="rules['activities']"
                  item-text="name"
                  return-object
                  multiple
                >
                  <template slot="selection" slot-scope="data">
                    <v-chip
                      :input-value="data.selected"
                      close
                      color="primary"
                      @input="removeActivity(data.item)"
                    >
                      <v-subheader class="body-2">
                        {{
                          $t('allActivities.' + `${data.item.name}`)
                            | capitalize
                        }}
                      </v-subheader>
                    </v-chip>
                  </template>
                </v-autocomplete>
                <v-slider
                  v-model="couch_number"
                  :rules="rules['couch_number']"
                  color="primary"
                  :label="$t('p.ped.couch_number')"
                  min="0"
                  max="15"
                  thumb-label
                />
              </v-flex>

              <v-flex
                :class="{
                  'sm6 pl-1 pb-0': $vuetify.breakpoint.smAndUp,
                  'xs12 pa-2': $vuetify.breakpoint.xsOnly
                }"
              >
                <v-card>
                  <app-map
                    ref="map"
                    v-model="markers"
                    :error-color="`info`"
                    :show="showEditPageDialog"
                    :help="$t('p.ped.map_help')"
                    :height="`40vh`"
                    :width="`100%`"
                    :center="mapConfig.center"
                    :max="mapConfig.markers.max"
                    :min="mapConfig.markers.min"
                    :errors="mapConfig.markers.errors"
                    @input="addMarker"
                  />
                </v-card>
              </v-flex>
              <v-flex v-if="images" xs12 mt-0 pt-0>
                <upload-button
                  :label="$t('p.ped.upload_image_hint') | capitalize"
                  :multiple="false"
                  :accept="rules.imageFormatsAllowed"
                  title="Browser"
                  :disabled="images.length > 5"
                  :errors="imagesErrors"
                  @formData="addImage"
                />
                <v-layout v-if="images.length > 0" row wrap>
                  <v-flex
                    v-for="(i, idx) in images"
                    :key="idx"
                    :class="{
                      'sm2 px-2': $vuetify.breakpoint.smAndUp,
                      'xs6 px-1 py-2': $vuetify.breakpoint.xsOnly
                    }"
                  >
                    <v-card class="rounded">
                      <v-img
                        :src="i.url"
                        :lazy-src="i.url"
                        :alt="i.alt"
                        aspect-ratio="1"
                        class="grey lighten-2"
                      >
                        <v-layout
                          slot="placeholder"
                          fill-height
                          align-center
                          justify-center
                          ma-0
                        >
                          <v-progress-circular
                            indeterminate
                            color="grey lighten-5"
                          />
                        </v-layout>
                        <v-text-field
                          :placeholder="$t('p.ped.image_alt_ph')"
                          :value="i.alt"
                          height="5"
                          solo
                          single-line
                          append-icon="mdi-close"
                          hide-details
                          color="grey"
                          background-color="rgba(255,255,255,0.7)"
                          :rules="rules['alt']"
                          @change="setImageAlt(idx, $event)"
                          @click:append="deleteImage(idx)"
                          @keypress.enter="validate"
                        />
                      </v-img>
                    </v-card>
                  </v-flex>
                </v-layout>
              </v-flex>
            </v-layout>
          </v-container>
        </v-form>
      </template>
    </v-card>
  </v-dialog>
</template>

<script>
  import UploadButton from 'components/utils/UploadButton'
  import AppMap from 'components/utils/AppMap'

  import {
    MODIFY_PAGE,
    PAGE_ADD_IMAGE,
    MODIFY_IMAGE_ALT,
    PAGE_DELETE_IMAGE,
    SAVE_PAGE,
    CANCEL_EDIT_PAGE,
    REMOVE_ACTIVITY
  } from 'store/pages/actions'

  import { mapMutations, mapActions, mapState } from 'vuex'
  import mapStatesTwoWay from 'plugins/mapStatesTwoWay'

  const NAMESPACE = 'pages/'

  export default {
    name: 'ProfilePageEditionDialog',
    components: { UploadButton, AppMap },
    props: { state: { type: String, default: 'edit' } },
    data() {
      return {
        isEditing: this.state === 'edit',

        showEditPageDialog: false,

        imagesErrors: [],

        markers: [],
        mapConfig: {
          zoom: 1,
          center: [46, -1],
          markers: {
            max: 1,
            min: 1,
            errors: {
              too_much: this.$t('p.ped.markers.too_much'),
              too_few: this.$t('p.ped.markers.too_few'),
              invalid: this.$t('p.ped.markers.invalid')
            }
          }
        },

        maxActivitiesAllowed: 3,
        errors: [],
        rules: {
          valid: true,
          invalidLocation: false,
          imageFormatsAllowed: 'image/jpeg, image/jpg, image/png, image/gif',
          imageSize: { w: 255, h: 255 },
          name: [
            (v) => !!v || this.$t('message.required', ['', this.$t('name')]),
            (v) =>
              (!!v && v.length > 6) || this.$t('message.length_above', [6]),
            (v) =>
              (!!v && v.length <= 40) || this.$t('message.length_below', [40]),
            (v) =>
              (!!v &&
                /^[a-zA-Z àáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð]+$/.test(
                  v
                )) ||
              this.$t('message.valid_chars_hint', ['a-zA-Z'])
          ],
          description: [
            (v) =>
              !!v || this.$t('message.required', ['e', this.$t('description')]),
            (v) =>
              /^[0-9a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,!?.'-]{0,75}$/i.test(
                v
              ) ||
              this.$t('message.invalid', [
                this.$t('the_f') + ' ' + this.$t('description')
              ])
          ],
          long_description: [
            (v) =>
              /^[0-9a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,!?.'-]{0,512}$/i.test(
                v
              ) || this.$t('message.invalid', ['description'])
          ],
          couch_number: [(v) => v < 15 || this.$t('p.ped.wow')],
          activities: [
            (v) =>
              !!v || this.$t('message.required', ['e', this.$t('activity')]),
            (v) =>
              v.length > 0 ||
              this.$t('message.required', ['e', this.$t('activity')]),
            (v) =>
              v.length <= this.maxActivitiesAllowed ||
              this.$t('message.too_much', [
                this.maxActivitiesAllowed,
                this.$t('activities')
              ])
          ],
          alt: [
            (v) =>
              /^[a-zA-Z0-9!? ]{0,15}$/.test(v) ||
              this.$t('message.valid_chars_hint', ['a-zA-Z0-9!? '])
            // v => !!v && v.length < 15 || this.$t("p.ped.invalid_image_alt")
          ]
        }
      }
    },
    computed: {
      ...mapState({ allActivities: (state) => state.profile.activities }),
      ...mapState({ edited_page: (state) => state.profile.pages.edited_page }),
      ...mapStatesTwoWay(
        'pages',
        {
          name: (state) => state.edited_page.name,
          description: (state) => state.edited_page.description,
          long_description: (state) => state.edited_page.long_description,
          lat: (state) => state.edited_page.lat,
          lng: (state) => state.edited_page.lng,
          couch_number: (state) => state.edited_page.couch_number,
          public: (state) => state.edited_page.public,
          activities: (state) => state.edited_page.activities,
          images: (state) => state.edited_page.images
        },
        function (payload) {
          this.MODIFY_PAGE(payload)
        }
      )
    },
    watch: {
      showEditPageDialog(v) {
        if (!v) return
        if (this.state === 'edit' && !!this.lat && !!this.lng) {
          this.showEditPageDialog = true
          this.markers = [{ lat: this.lat, lng: this.lng }]
          this.mapConfig.zoom = 5
          this.mapConfig.center = [this.lat, this.lng]
        }
      }
    },
    mounted() {
      if (navigator.geolocation) {
        var self = this
        navigator.geolocation.getCurrentPosition(function (position) {
          self.mapConfig.center = {
            lat: position.coords.latitude,
            lng: position.coords.longitude
          }
        })
      }
    },
    methods: {
      ...mapMutations(NAMESPACE, [
        MODIFY_IMAGE_ALT,
        MODIFY_PAGE,
        CANCEL_EDIT_PAGE,
        PAGE_ADD_IMAGE,
        REMOVE_ACTIVITY
      ]),
      ...mapActions(NAMESPACE, [SAVE_PAGE, PAGE_DELETE_IMAGE]),
      validate() {
        this.rules.invalidLocation = false
        if (!this.$refs.form.validate()) {
          return
        }

        if (this.$refs.map.error) return

        if (
          this.markers.length > this.mapConfig.markers.max ||
          this.markers.length < this.mapConfig.markers.min
        ) {
          return
        }

        if (this.images.length === 0) {
          return (this.imagesErrors = [
            this.$t('message.required', ['e', this.$t('image')])
          ])
        }

        this.submit()
      },
      submit() {
        this.$loader(true)
        this.SAVE_PAGE(this.state)
          .then(() => {
            this.$loader(false)
            this.showEditPageDialog = false
            this.$emit('page-saved', true)
            this.delMarker()
          })
          .catch(() => {
            this.$loader(false)
            this.$emit('page-saved', false)
          })
      },
      removeActivity(activity) {
        this.REMOVE_ACTIVITY(activity)
      },
      addMarker(markers) {
        if (markers.length === 0) return this.delMarker()
        this.MODIFY_PAGE({ lat: markers[0].lat })
        this.MODIFY_PAGE({ lng: markers[0].lng })
      },
      delMarker() {
        this.MODIFY_PAGE({ lat: null })
        this.MODIFY_PAGE({ lng: null })
      },
      clear() {
        this.$refs.form.reset()
      },
      setImageAlt(idx, value) {
        this.MODIFY_IMAGE_ALT({
          idx: idx,
          value: value
        })
      },
      cancelEdit() {
        this.CANCEL_EDIT_PAGE()
        this.showEditPageDialog = false
        this.imagesErrors = []
        this.markers = []
        this.zoom = 5
      },
      addImage(formData) {
        if (this.images.length > 5) {
          this.$snackbar(this.$t('p.pde.max_images'))
          return
        }

        var file = formData.get('file')
        if (file instanceof File) {
          if (file.size > 500000) {
            this.imagesErrors.push(
              this.$t('message.too_big', [this.$t('image'), '500ko'])
            )
            return
          }

          var exists = this.images.filter(
            (i) =>
              (i.url && i.url.indexOf(file.name) > -1) ||
              (i.File && i.File.indexOf(file.name) > -1)
          ).length
          if (exists > 0) {
            this.imagesErrors.push(this.$t('message.exist', ['image']))
            return
          }

          var that = this
          var reader = new FileReader()
          reader.onload = function (e) {
            that.PAGE_ADD_IMAGE({
              url: e.target.result,
              File: file.name
            })
            that.imagesErrors = []
          }

          reader.readAsDataURL(file)
        }
      },
      deleteImage(idx) {
        this.PAGE_DELETE_IMAGE(idx)
          .then(() => {
            this.$snackbar(
              this.$t('message.success_deleting', [this.$t('image')])
            )
          })
          .catch(() => {
            this.$snackbar(
              this.$t('message.error_deleting', [this.$t('image')])
            )
          })
      }
    }
  }
</script>

<style lang="scss">
  .rounded {
    border-radius: 10px;
  }

  .image-alt-in {
    position: absolute;
    line-height: 27px;
    background-color: rgba(#fff, 0.8);
    width: 100%;
    bottom: 0;
    padding: 8px;

    &:focus {
      outline: none;
    }
  }
</style>
