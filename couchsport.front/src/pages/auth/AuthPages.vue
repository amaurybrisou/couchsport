<template>
  <v-container fluid fill-height>
    <v-layout justify-center align-center>
      <v-flex xs12 sm8 md4>
        <auth-form
          v-bind="{
            title,
            buttonMessage,
            welcome,
            submit,
            flat,
            color,
            errors
          }"
        />
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import { AUTH_REQUEST, AUTH_SIGNUP } from 'store/auth/actions'
  import AuthForm from 'components/auth/AuthForm'

  export default {
    name: 'AuthPages',
    components: { AuthForm },
    props: {
      welcome: { type: String, default: '' }
    },
    data() {
      return { flat: false, color: 'primary', errors: [] }
    },
    computed: {
      submit() {
        return this.$route.name === 'signup' ? this.SignUp : this.Login
      },
      title() {
        return this.$t(this.$route.name)
      },
      buttonMessage() {
        return this.$t(this.$route.name)
      }
    },
    methods: {
      SignUp: function (user) {
        this.errors = []
        this.$store
          .dispatch(AUTH_SIGNUP, user)
          .then(() => {
            this.$router.push({
              name: 'login',
              params: {
                welcome: this.$t('message.signup_success_welcome')
              }
            })
          })
          .catch((error) => {
            this.errors = [error]
          })
      },
      Login: function (user) {
        this.errors = []
        this.$store
          .dispatch(AUTH_REQUEST, user)
          .then(() => {
            this.$router.push({ name: 'informations' })
          })
          .catch((error) => {
            this.errors = [error]
          })
      }
    }
  }
</script>
