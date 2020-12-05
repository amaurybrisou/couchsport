<template>
  <nav>
    <v-system-bar color="primary" window height="36vh">
      <v-btn :to="{ name: 'home' }" exact text tile>
        <v-icon>mdi-home</v-icon>
      </v-btn>
      <v-btn :to="{ name: 'explore' }" text tile>
        {{ $t('explore') }}
      </v-btn>
      <language-button />

      <v-spacer />

      <a
        v-if="
          isProfileLoaded && unread_message && $route.hash !== '#conversations'
        "
        class="new-message-link"
        @click="goToConversations"
      >
        <v-chip small color="info" class="mx-3">
          {{ unread_message }}
          <v-icon class="ml-3">mdi-mail</v-icon>
        </v-chip>
      </a>
      <v-menu
        v-if="isAuthenticated && isProfileLoaded"
        v-ws.connect="`${getProfile.id}`"
        open-on-hover
        offset-y
        left
        transition="slide-y-transition"
        style="z-index: 500"
      >
        <template v-slot:activator="{ on }">
          <v-btn icon v-on="on">
            <v-icon>mdi-account</v-icon>
          </v-btn>
        </template>
        <v-list style="cursor: pointer">
          <v-list-item-group>
            <v-list-item
              v-for="link in links"
              :key="link.name"
              :to="{ name: 'profile', hash: link.to }"
            >
              {{ $t(`${link.name}`) | capitalize }}
            </v-list-item>

            <v-list-item :to="{ name: 'about' }">
              {{ $t('about') | capitalize }}
            </v-list-item>

            <v-list-item @click="logout">
              {{ $t('logout') | capitalize }}
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-menu>

      <v-btn :to="{ name: 'about' }" text tile>
        {{ $t('about') }}
      </v-btn>

      <template v-if="!isAuthenticated">
        <v-btn :to="{ name: 'signup' }" text tile>
          {{ $t('signup') | capitalize }}
        </v-btn>
        <v-btn v-if="!authLoading" :to="{ name: 'login' }" text tile>
          {{ $t('login') | capitalize }}
        </v-btn>
      </template>
    </v-system-bar>
  </nav>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  import { AUTH_LOGOUT } from 'actions/auth'
  import LanguageButton from 'components/nav/LanguageButton'

  export default {
    name: 'AppNav',
    components: { LanguageButton },
    data() {
      return {
        links: [
          { auth: true, to: '#informations', name: 'profile' },
          { auth: true, to: '#activities', name: 'activities' },
          { auth: true, to: '#conversations', name: 'conversations' },
          { auth: true, to: '#pages', name: 'pages' }
        ]
      }
    },
    computed: {
      ...mapGetters(['getProfile', 'isAuthenticated', 'isProfileLoaded']),
      ...mapState({
        unread_message: (state) => state.profile.conversations.unread,
        authLoading: (state) => state.auth.status === 'loading',
        name: (state) =>
          `${state.user.profile.Firstname} ${state.user.profile.Lastname}`
      })
    },

    methods: {
      goToConversations() {
        this.$router.push({ name: 'profile', hash: '#conversations' })
      },

      logout: function () {
        this.$store
          .dispatch(AUTH_LOGOUT)
          .then(() => this.$router.push({ name: 'home' }))
      }
    }
  }
</script>

<style lang="scss">
  .new-message-link {
    text-decoration: none;
    color: none;
    outline: none;

    span:focus:after {
      background: none !important;
    }
  }
</style>
