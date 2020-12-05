<template>
  <v-form @submit="submitForm">
    <v-card>
      <v-toolbar :color="color">
        <v-toolbar-title>{{ title | capitalize }}</v-toolbar-title>
      </v-toolbar>
      <v-card-text>
        <div v-if="errors.length" color="error">
          <v-alert v-for="(err, i) in errors" :key="i" type="error">
            {{ $t(err) }}
          </v-alert>
        </div>
        <div v-if="welcome.length > 0" color="info">
          <v-alert type="info">
            {{ welcome | capitalize }}
          </v-alert>
        </div>

        <v-form ref="form" v-model="valid" @keypress.enter.native="submitForm">
          <v-text-field
            v-show="!email"
            v-model="user.email"
            :label="$t('email') | capitalize"
            type="text"
            name="email"
            :rules="emailRules"
            autocomplete="email"
          />
          <v-text-field
            v-model="user.password"
            :label="$t('password') | capitalize"
            :type="'password'"
            name="password"
            counter="8"
            :rules="passwordRules"
            autocomplete="current-password"
          />
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          name="submit"
          :color="color"
          :text="flat"
          :disabled="!valid"
          @click="submitForm"
        >
          {{ buttonMessage | capitalize }}
        </v-btn>
        <v-btn
          v-show="email"
          name="cancel"
          :color="color"
          :text="flat"
          @click="$emit('hide-change-password-dialog')"
        >
          {{ $t('cancel') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  export default {
    name: 'AuthForm',
    props: {
      welcome: { type: String, default: '', required: false },
      title: { type: String, default: 'signup' },
      buttonMessage: { type: String, default: 'signup' },
      flat: { type: Boolean, default: false },
      color: { type: String, default: 'primary' },
      submit: {
        type: Function,
        default: () => {
          throw 'Not Implemented'
        }
      }
    },
    data() {
      return {
        errors: [],
        user: { email: '', password: '' },
        valid: false,
        emailRules: [
          (v) => !!v || this.$t('message.required', ['', this.$t('email')]),
          (v) =>
            /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) ||
            this.$t('message.invalid', [this.$t('email')])
        ],
        passwordRules: [
          (v) => !!v || this.$t('message.required', ['', this.$t('password')]),
          (v) =>
            /^(?=.*\d)(?=.*[_!?,]?)(?=.*[a-z])(?=.*[A-Z])[0-9a-zA-Z_?,!]{8,}$/.test(
              v
            ) || this.$t('message.password_hint', [8])
        ]
      }
    },
    computed: {
      ...mapGetters(['isAuthenticated', 'isProfileLoaded']),
      ...mapState({
        email: (state) => state.auth.email
      })
    },
    methods: {
      submitForm() {
        this.submit(this.user)
        this.$emit('submit', this.user)
      }
    }
  }
</script>
