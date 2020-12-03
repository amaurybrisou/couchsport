<template>
  <v-container fluid fill-height>
    <v-layout justify-center align-center>
      <v-flex xs12 sm8 md4>
        <auth-form
          :title="$t('signup') | capitalize"
          :button-message="$t('signup') | capitalize"
          :errors="errors"
          @submit="submit"
        />
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import { repositoryFactory } from 'repos/repositoryFactory'

  import AuthForm from './AuthForm'
  const userRepository = repositoryFactory.get('user')

  export default {
    name: 'SignUp',
    components: { AuthForm },
    data() {
      return { errors: [] }
    },
    methods: {
      submit(user) {
        userRepository
          .create(user)
          .then(() => {
            this.$router.push({
              name: 'login',
              params: { welcome: this.$t('message.signup_success_welcome') }
            })
          })
          .catch(({ response: { data } }) => {
            this.errors = []
            this.errors.push(data)
          })
      }
    }
  }
</script>
