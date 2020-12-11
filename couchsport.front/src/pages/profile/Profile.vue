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
        >
          <v-tab v-for="tab of tabs" :key="tab.id" :to="{ name: tab.name }">
            {{ tab.text }}
          </v-tab>
          <v-tabs-items>
            <keep-alive><router-view /></keep-alive>
          </v-tabs-items>
          <!-- <v-tab to="informations" class="subheader">
            {{ $t('personal_informations') }}
          </v-tab>
          <v-tab to="#activities" class="subheader">
            {{ $t('activities') }}
          </v-tab>
          <v-tab
            to="#conversations"
            class="subheader"
            @click="$messenger.setMessagesRead()"
          >
            {{ $t('conversations') }}
          </v-tab>
          <v-tab to="#pages" class="subheader">
            {{ $t('pages') }}
          </v-tab> -->

          <!-- <v-tabs-items v-model="activeTab">
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
          </v-tabs-items> -->
        </v-tabs>
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
  // import Informations from 'pages/profile/Informations'
  // import Activities from 'pages/profile/Activities'
  // import Pages from 'pages/profile/Pages'
  // import Conversations from 'pages/profile/Conversations'
  import { mapGetters } from 'vuex'
  import {
    SET_ACTIVITIES,
    SET_LANGUAGES,
    PROFILE_REQUEST
  } from 'store/profile/actions'

  export default {
    name: 'Profile',
    // components: {
    //   Informations,
    //   Activities,
    //   Pages,
    //   Conversations
    // },
    data() {
      return {
        activeTab: 'informations',
        tabs: [
          { name: 'informations', text: this.$t('personal_informations') },
          { name: 'activities', text: this.$t('activities') },
          { name: 'conversations', text: this.$t('conversations') },
          { name: 'pages', text: this.$t('pages') }
        ]
      }
    },
    computed: {
      ...mapGetters(['isProfileLoaded', 'getProfile'])
    },
    created() {
      Promise.all([
        this.$store.dispatch(PROFILE_REQUEST),
        this.$store.dispatch(
          SET_ACTIVITIES,
          this.$route.params.locale,
          this.$route.params.locale
        ),
        this.$store.dispatch(
          SET_LANGUAGES,
          this.$route.params.locale,
          this.$route.params.locale
        )
      ]).catch(() =>
        this.$route.push({ name: 'home', locale: this.getProfile.locale })
      )
    }
  }
</script>
