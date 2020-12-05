<template>
  <v-container fill-height>
    <v-layout justify-center align-center>
      <v-flex xs8 sm5 md6>
        <v-card>
          <v-toolbar flat color="primary">
            <h1 class="display-2 font-weight-thin">CoSport</h1>
          </v-toolbar>
          <v-card-text>
            {{ $t('p.about.content') }}
            <v-layout row wrap>
              <v-spacer />
              <v-btn text :to="{ name: 'signup' }">
                {{ $t('signup') }}
              </v-btn>
            </v-layout>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  export default {
    name: 'About'
  }
</script>
