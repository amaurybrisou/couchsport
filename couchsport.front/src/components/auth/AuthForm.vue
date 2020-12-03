<template>
  <v-form @submit="submit">
    <v-card>
      <v-toolbar :color="color">
        <v-toolbar-title>{{ title }}</v-toolbar-title>
      </v-toolbar>
      <v-card-text>
        <div v-if="errors.length" color="error">
          <v-alert v-for="(err, i) in errors" :key="i" type="error">
            {{ err }}
          </v-alert>
        </div>
        <div v-if="welcome" color="info">
          <v-alert type="info">
            {{ welcome | capitalize }}
          </v-alert>
        </div>

        <v-form ref="form" v-model="valid" @keypress.enter.native="submit">
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
        <v-btn :color="color" :text="flat" :disabled="!valid" @click="submit">
          {{ buttonMessage }}
        </v-btn>
        <v-btn
          v-show="email"
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
  export default {
    name: 'AuthForm',
    props: {
      welcome: { type: String, default: null, required: false },
      title: { type: String, default: 'login' },
      buttonMessage: { type: String, default: 'login' },
      errors: { type: Array, default: () => [], required: false },
      flat: { type: Boolean, default: false },
      color: { type: String, default: 'primary' }
    },
    data() {
      return {
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
      user() {
        return { email: this.email, password: '' }
      },
      email() {
        return this.$store.state.auth.email || ''
      },
      emailErrors() {
        const errors = []
        if (!this.$v.user.email.$dirty) return errors
        !this.$v.user.email.email &&
          errors.push(this.$t('message.invalid', [this.$t('email')]))
        !this.$v.user.email.required &&
          errors.push(this.$t('message.required', ['', this.$t('email')]))
        return errors
      },
      passwordErrors() {
        const errors = []
        if (!this.$v.user.password.$dirty) return errors
        !this.$v.user.password.maxLength &&
          errors.push(this.$t('message.password_hint', [8]))
        !this.$v.user.password.required &&
          errors.push(this.$t('message.required', ['', this.$t('password')]))
        return errors
      }
    },
    methods: {
      submit() {
        this.$emit('submit', this.user)
      }
    }
  }
</script>
