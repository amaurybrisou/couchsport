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
  import { mapGetters } from 'vuex'

  export default {
    components: { AuthForm },
    props: {
      welcome: { type: String, default: '' }
    },
    data() {
      return { flat: false, color: 'primary' }
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
      },
      ...mapGetters(['errors'])
    },
    methods: {
      SignUp: async function (user) {
        const response = await this.$store.dispatch(AUTH_SIGNUP, user)

        if (!response) return

        this.$router.push({
          name: 'login',
          params: {
            welcome: this.$t('message.signup_success_welcome')
          }
        })
      },
      Login: async function (user) {
        const response = await this.$store.dispatch(AUTH_REQUEST, user)
        if (!response) return
        this.$router.push({ name: 'profile' })
      }
    }
  }
</script>
