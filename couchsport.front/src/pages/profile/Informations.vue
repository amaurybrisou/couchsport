<template>
  <v-container>
    <v-form
      v-if="isProfileLoaded"
      ref="form"
      v-model="rules.valid"
      @keypress.enter.native="submit"
    >
      <v-layout row wrap align-center justify-center>
        <v-flex xs6 sm5 md4 xl3 mt-5 mr-5>
          <upload-button
            :accept="rules.imageFormatsAllowed"
            @formData="handleImage"
          >
            <v-card slot="appearance" text tile class="profile-avatar">
              <v-img
                :src="avatar"
                :alt="username"
                :value="avatar"
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
                  <!-- <v-progress-circular indeterminate color="grey lighten-5"></v-progress-circular> -->
                  <v-icon large>mdi-face-profile</v-icon>
                  <v-subheader color="warning">
                    {{ $t('message.required', ['e', $t('image')]) }}
                  </v-subheader>
                </v-layout>
              </v-img>
            </v-card>
          </upload-button>
          <v-input
            v-model="avatar"
            :rules="rules['avatar']"
            hidden
            class="right"
          />
        </v-flex>
        <v-flex xs11 sm5 md12 lg6 xl4>
          <v-text-field
            v-model="email"
            text
            disabled
            readonly
            :label="$t('fields.email') | capitalize"
          />
          <v-text-field
            v-model="username"
            text
            autofocus
            :rules="rules['username']"
            :label="$t('fields.username') | capitalize"
            @keypress.enter.native="submit"
          />
          <v-text-field
            v-model="firstname"
            text
            :rules="rules['firstname']"
            :label="$t('fields.firstname') | capitalize"
            @keypress.enter.native="submit"
          />
          <v-text-field
            v-model="lastname"
            text
            :rules="rules['lastname']"
            :label="$t('fields.lastname') | capitalize"
            @keypress.enter.native="submit"
          />
          <v-select
            v-model="gender"
            :items="[``, `Male`, `Female`]"
            :label="$t('fields.gender') | capitalize"
            @keypress.enter.native="submit"
          />
        </v-flex>
        <v-flex xs11 sm10 md10 xl4>
          <v-text-field
            v-model="street_name"
            text
            :label="$t('fields.streetname') | capitalize"
            :rules="rules['street_name']"
            @keypress.enter.native="submit"
          />
          <v-text-field
            v-model="city"
            text
            :label="$t('fields.city') | capitalize"
            :rules="rules['city']"
            @keypress.enter.native="submit"
          />
          <v-text-field
            v-model="zip_code"
            text
            :label="$t('fields.zipcode') | capitalize"
            :rules="rules['zip_code']"
            @keypress.enter.native="submit"
          />
          <v-text-field
            v-model="country"
            text
            :label="$t('fields.country') | capitalize"
            :rules="rules['country']"
            @keypress.enter.native="submit"
          />
          <v-text-field
            v-model="phone"
            text
            :rules="rules['phone']"
            :label="$t('fields.phone') | capitalize"
            @keypress.enter.native="submit"
          />
        </v-flex>
        <v-flex xs11 sm10 md10 xl11>
          <v-autocomplete
            v-model="languages"
            :items="allLanguages"
            :label="$t('languages') | capitalize"
            return-object
            item-text="native_name"
            multiple
          />
        </v-flex>
        <v-flex xs12 md12 mb-5>
          <v-btn
            block
            text
            color="warning"
            @click="showChangePasswordDialog = true"
          >
            {{ $t('change_password') }}
          </v-btn>
        </v-flex>
        <v-flex xs11 md12 mb-5>
          <v-btn block text color="success" @click="submit">
            {{ $t('save') }}
          </v-btn>
        </v-flex>
      </v-layout>
    </v-form>

    <!-- Warn section  used to display application state (saving and success) -->
    <v-dialog v-model="showChangePasswordDialog" persistent width="350">
      <auth-form
        :title="$t('change_password') | capitalize"
        :button-message="$t('change_password') | capitalize"
        text
        :color="`warning`"
        :submit="changePassword"
        @hide-change-password-dialog="showChangePasswordDialog = false"
      />
    </v-dialog>
  </v-container>
</template>

<script>
  import { mapGetters, mapState, mapMutations, mapActions } from 'vuex'

  import mapStatesTwoWay from 'plugins/mapStatesTwoWay'

  import { MODIFY_PROFILE, SAVE_PROFILE } from 'store/profile/actions'

  import { AUTH_CHANGE_PASSWORD } from 'store/auth/actions'

  import UploadButton from 'components/utils/UploadButton'
  import AuthForm from 'components/auth/AuthForm'

  export default {
    name: 'PersonalInformation',
    components: { UploadButton, AuthForm },
    data() {
      return {
        showChangePasswordDialog: false,

        rules: {
          valid: false,
          imageFormatsAllowed: 'image/jpeg, image:jpg, image/png, image/gif',
          username: [
            (v) =>
              (v.length >= 6 && v.length < 15) ||
              this.$t('message.length_between', [
                this.$t('fields.username'),
                6,
                15
              ]),
            (v) =>
              /^[àéèïîôoa-zA-Z]{6,15}$/.test(v) ||
              this.$t('message.valid_chars_hint', ['àéèïîôoa-zA-Z'])
          ],
          firstname: [
            (v) =>
              v.length < 35 ||
              this.$t('message.length_below', [
                this.$t('fields.firstname'),
                35
              ]),
            (v) =>
              /^[àéèêïîôo a-zA-Z]{0,35}$/.test(v) ||
              this.$t('message.valid_chars_hint', ['àéèêïîôo a-zA-Z'])
          ],
          lastname: [
            (v) =>
              v.length < 35 ||
              this.$t('message.length_below', [this.$t('fields.lastname'), 35]),
            (v) =>
              /^[àéèêïîôo a-zA-Z]{0,35}$/.test(v) ||
              this.$t('message.valid_chars_hint', ['àéèêïîôo a-zA-Z'])
          ],
          avatar: [
            (v) => !!v || this.$t('message.required', ['e', this.$t('image')]),
            (v) =>
              /(?:png|jpg|jpeg|gif)$/i.test(v) ||
              // eslint-disable-next-line no-useless-escape
              /^\s*data:([a-z]+\/[a-z]+(;[a-z\-]+\=[a-z\-]+)?)?(;base64)?,[a-z0-9\!\$\&\'\,\(\)\*\+\,\;\=\-\.\_\~\:\@\/\?\%\s]*\s*$/i.test(
                v
              ) ||
              this.$t('message.invalid', [this.$t('image_link')])
          ],
          zip_code: [
            (v) =>
              v.length < 35 ||
              this.$t('message.length_below', [this.$t('fields.zipcode'), 35]),
            (v) =>
              /^[0-9a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,!?.'-]{0,35}$/.test(
                v
              ) || this.$t('message.invalid', [this.$t('fields.zipcode')])
          ],
          street_name: [
            (v) =>
              v.length < 50 ||
              this.$t('message.length_below', [
                this.$t('fields.streetname'),
                50
              ]),
            (v) =>
              /^[0-9a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,!?.'-]{0,35}$/.test(
                v
              ) || this.$t('message.invalid', [this.$t('fields.streetname')])
          ],
          city: [
            (v) =>
              v.length < 35 ||
              this.$t('message.length_below', [this.$t('fields.city'), 35]),
            (v) =>
              /^[0-9a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,!?.'-]{0,35}$/.test(
                v
              ) || this.$t('message.invalid', [this.$t('fields.city')])
          ],
          country: [
            (v) =>
              v.length < 35 ||
              this.$t('message.length_below', [this.$t('fields.country'), 35]),
            (v) =>
              /^[0-9a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,!?.'-]{0,35}$/.test(
                v
              ) || this.$t('message.invalid', [this.$t('fields.country')])
          ],
          phone: [
            (v) =>
              v.length < 35 ||
              this.$t('message.length_below', [this.$t('fields.phone'), 35]),
            (v) =>
              /^[0-9a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,!?.'-]{0,35}$/.test(
                v
              ) || this.$t('message.invalid', [this.$t('fields.phone')])
          ]
        }
      }
    },
    computed: {
      ...mapStatesTwoWay(
        null,
        {
          username: (state) => state.profile.profile.username,
          firstname: (state) => state.profile.profile.firstname,
          lastname: (state) => state.profile.profile.lastname,
          gender: (state) => state.profile.profile.gender,
          street_name: (state) => state.profile.profile.street_name,
          city: (state) => state.profile.profile.city,
          zip_code: (state) => state.profile.profile.zip_code,
          country: (state) => state.profile.profile.country,
          phone: (state) => state.profile.profile.phone,
          languages: (state) => state.profile.profile.languages
        },
        function (payload) {
          this.MODIFY_PROFILE(payload)
        }
      ),
      ...mapGetters({
        isProfileLoaded: 'isProfileLoaded'
      }),
      ...mapState({
        email: (state) => state.auth.email,
        allLanguages: (state) => state.profile.languages
      }),
      avatar: {
        get() {
          return this.$store.state.profile.profile.avatar
        }
      }
    },
    methods: {
      ...mapActions([SAVE_PROFILE, AUTH_CHANGE_PASSWORD]),
      ...mapMutations([MODIFY_PROFILE]),
      submit() {
        if (!this.$refs.form.validate()) {
          return
        }
        this.$loader(true)
        this.SAVE_PROFILE()
          .then(() => {
            this.$snackbar(this.$t('message.success_saving', ['profile']))
          })
          .catch(() => {
            this.$snackbar(this.$t('message.error_saving', ['profile']))
          })
          .finally(() => {
            this.$loader(false)
          })
      },
      changePassword(user) {
        this.$loader(true)
        this[AUTH_CHANGE_PASSWORD](user)
          .then(() => {
            this.$snackbar(
              this.$t('message.success_updating', [this.$t('password')])
            )
            this.showChangePasswordDialog = false
          })
          .catch(() => {
            this.$snackbar(
              this.$t('message.error_updating', [this.$t('password')])
            )
          })
          .finally(() => {
            this.$loader(false)
          })
      },
      handleImage(formData) {
        var file = formData.get('file')
        if (file instanceof File && file.size) {
          if (file.size > 100000) {
            this.$snackbar(this.$t('message.too_big', ['image', '100ko']))
            return
          }
          var that = this
          var reader = new FileReader()
          reader.onload = function (error) {
            that.$loader(false)
            that.MODIFY_PROFILE({ avatar: error.target.result })
            that.MODIFY_PROFILE({ avatar_file: file.name })
          }

          reader.onerror = function (error) {
            console.error(error)
            that.$snackbar('message.error_reading', ['image'])
          }

          that.$loader(true)
          reader.readAsDataURL(file)
        }
      }
    }
  }
</script>

<style lang="scss">
  .profile-avatar {
    border-radius: 50%;
  }
</style>
