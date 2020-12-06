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
  export default {
    props: {
      mutator: {
        type: Function,
        default: () => {
          throw 'Not Implemented'
        }
      }
    },
    data() {
      return {
        languages: {
          en: { title: 'English' },
          fr: { title: 'Français' }
        }
      }
    },
    mounted() {
      this.mutator({ locale: this.$i18n.locale })
    },
    methods: {
      changeLocale(locale) {
        this.$i18n.locale = locale
        this.mutator({ locale: locale })

        this.$router.push({
          name: this.$route.name,
          params: { locale: locale }
        })
      }
    }
  }
</script>
