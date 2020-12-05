<template>
  <v-menu offset-y left transition="slide-y-transition" style="z-index: 500">
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
</template>

<script>
  import { MODIFY_PROFILE } from 'actions/profile'
  import { mapMutations } from 'vuex'

  export default {
    data() {
      return {
        languages: {
          en: { title: 'English' },
          fr: { title: 'Français' }
        }
      }
    },
    created() {
      this.MODIFY_PROFILE({ locale: this.$i18n.locale })
    },
    methods: {
      ...mapMutations([MODIFY_PROFILE]),
      changeLocale(locale) {
        this.$i18n.locale = locale
        this.MODIFY_PROFILE({ locale: locale })

        this.$router.push({
          name: this.$route.name,
          params: { locale: locale }
        })
      }
    }
  }
</script>
