<template>
  <nav>
    <v-system-bar color="primary" window height="36vh">
      <v-btn :to="{ name: 'home' }" exact text tile>
        <v-icon>mdi-home</v-icon>
      </v-btn>
      <v-btn :to="{ name: 'explore' }" text tile>
        {{ $t('explore') }}
      </v-btn>
      <v-menu
        offset-y
        left
        transition="slide-y-transition"
        style="z-index: 500"
      >
        <template v-slot:activator="{ on: menu }">
          <v-btn text color="primary" dark v-on="{ ...menu }">
            <v-icon class="mr-1">mdi-web</v-icon>
            {{ $i18n.locale.toUpperCase() }}
          </v-btn>
        </template>
        <v-list>
          <v-list-item-group>
            <v-list-item
              v-for="(item, i) in languages"
              :key="i"
              @click="changeLocale(i)"
            >
              <v-list-item-content>{{ item.title }}</v-list-item-content>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-menu>
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
        v-if="isProfileLoaded"
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

            <v-list-item
              v-if="isAuthenticated && isProfileLoaded"
              @click="logout"
            >
              {{ $t('logout') | capitalize }}
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-menu>

      <v-btn
        v-if="!isAuthenticated && !isProfileLoaded"
        :to="{ name: 'about' }"
        text
        tile
      >
        {{ $t('about') }}
      </v-btn>
      <v-btn v-if="!isAuthenticated" :to="{ name: 'signup' }" text tile>
        {{ $t('signup') | capitalize }}
      </v-btn>
      <v-btn
        v-if="!isAuthenticated && !authLoading"
        :to="{ name: 'login' }"
        text
        tile
      >
        {{ $t('login') | capitalize }}
      </v-btn>

      <!-- <v-menu
        v-else
        open-on-hover
        offset-y
        left
        transition="slide-y-transition"
        style="z-index: 500"
      >
        <template v-slot:activator="{ on }">
          <v-btn icon v-on="on">
            <v-icon>account_box</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item class="v-list__item--link" :to="{ name: 'about' }">
            <v-list-item-title>{{
              $t('about') | capitalize
            }}</v-list-item-title>
          </v-list-item>
          <v-list-item
            v-if="!isAuthenticated"
            :to="{ name: 'login' }"
            class="v-list__item--link"
          >
            <v-list-item-title>{{
              $t('login') | capitalize
            }}</v-list-item-title>
          </v-list-item>
          <v-list-item
            v-if="!isAuthenticated"
            :to="{ name: 'signup' }"
            class="v-list__item--link"
          >
            <v-list-item-title>{{
              $t('signup') | capitalize
            }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu> -->
    </v-system-bar>
  </nav>
</template>

<script>
  import { mapGetters, mapState, mapMutations } from 'vuex'
  import { AUTH_LOGOUT } from 'actions/auth'
  import { MODIFY_PROFILE } from 'actions/profile'

  export default {
    name: 'AppNav',
    data() {
      return {
        links: [
          { auth: true, to: '#informations', name: 'profile' },
          { auth: true, to: '#activities', name: 'activities' },
          { auth: true, to: '#conversations', name: 'conversations' },
          { auth: true, to: '#pages', name: 'pages' }
        ],
        languages: {
          en: { title: 'English' },
          fr: { title: 'Français' }
        }
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
    mounted() {
      this.MODIFY_PROFILE({ locale: this.$i18n.locale })
    },
    methods: {
      ...mapMutations([MODIFY_PROFILE]),
      goToConversations() {
        this.$messenger.setMessagesRead()
        this.$router.push({ name: 'profile', hash: '#conversations' })
      },
      changeLocale(locale) {
        this.$i18n.locale = locale
        this.MODIFY_PROFILE({ locale: locale })

        this.$router.push({
          name: this.$route.name,
          params: { locale: locale }
        })
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
