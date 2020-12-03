<template>
  <v-container fluid fill-height>
    <v-layout justify-center align-center>
      <v-flex xs12 sm8 md4>
        <auth-form
          :title="$t('login') | capitalize"
          :button-message="$t('login') | capitalize"
          :welcome="welcome"
          :errors="errors"
          @submit="submit"
        />
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import { AUTH_REQUEST } from 'actions/auth'
  import { mapActions } from 'vuex'
  import AuthForm from './AuthForm'

  export default {
    name: 'Login',
    components: { AuthForm },
    props: { welcome: { type: String, default: '' } },
    data() {
      return { errors: [] }
    },
    methods: {
      ...mapActions([AUTH_REQUEST]),
      submit(user) {
        this.AUTH_REQUEST(user)
          .then(() => {
            this.$router.push({ name: 'profile' })
          })
          .catch((data) => {
            console.log(data)
            this.errors = []
            this.errors.push(data)
          })
      }
    }
  }
</script>
