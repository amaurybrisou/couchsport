<template>
  <div v-if="isProfileLoaded" id="profile" fill-height>
    <v-layout justify-center fill-height>
      <v-flex xs12>
        <v-tabs
          slot="extension"
          v-model="activeTab"
          slider-color="primary"
          color="primary"
          centered
          show-arrows
        >
          <v-tab to="#informations" href="#informations" class="subheader">
            {{ $t('personal_informations') }}
          </v-tab>
          <v-tab to="#activities" href="#activities" class="subheader">
            {{ $t('activities') }}
          </v-tab>
          <v-tab
            to="#conversations"
            href="#conversations"
            class="subheader"
            @click="$messenger.setMessagesRead()"
          >
            {{ $t('conversations') }}
          </v-tab>
          <v-tab to="#pages" href="#pages" class="subheader">
            {{ $t('pages') }}
          </v-tab>
          <v-tabs-items v-model="activeTab">
            <v-tab-item value="informations">
              <informations />
            </v-tab-item>
            <v-tab-item value="activities">
              <activities />
            </v-tab-item>
            <v-tab-item value="conversations">
              <conversations />
            </v-tab-item>
            <v-tab-item value="pages">
              <pages />
            </v-tab-item>
          </v-tabs-items>
        </v-tabs>
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
  import Informations from 'pages/profile/Informations'
  import Activities from 'pages/profile/Activities'
  import Pages from 'pages/profile/Pages'
  import Conversations from 'pages/profile/Conversations'
  import { mapGetters, mapActions } from 'vuex'
  import { GET_ACTIVITIES, GET_LANGUAGES } from 'actions/profile'

  export default {
    name: 'Profile',
    components: {
      Informations,
      Activities,
      Pages,
      Conversations
    },
    data() {
      return {
        activeTab: 'informations'
      }
    },
    computed: {
      ...mapGetters(['isProfileLoaded', 'getProfile'])
    },
    mounted() {
      this.GET_ACTIVITIES(this.$route.params.locale, this.$route.params.locale)
      this.GET_LANGUAGES(this.$route.params.locale, this.$route.params.locale)
    },
    methods: {
      ...mapActions([GET_ACTIVITIES, GET_LANGUAGES])
    }
  }
</script>
