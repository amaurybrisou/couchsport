<template>
  <v-container fluid fill-height>
    <v-layout justify-center align-center>
      <v-flex xs12 sm8 md4>
        <login v-if="$route.name === 'login'" v-bind="$attrs" />
        <sign-up v-if="$route.name === 'signup'" v-bind="$attrs" />
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import Login from 'components/auth/Login'
  import SignUp from 'components/auth/SignUp'
  export default {
    components: { Login, SignUp }
  }
</script>
