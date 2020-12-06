<template>
  <v-container fluild grid-list-xs>
    <v-layout pb-3>
      <div>{{ $t('p.activities.hint') | capitalize }}</div>
    </v-layout>
    <v-divider />
    <v-layout v-if="activities" wrap>
      <v-flex v-for="(item, i) in activities" :key="i" xs6 md2 px-5>
        <v-checkbox
          v-model="selected_activities"
          height="0"
          :label="$t(`allActivities.${item.name}`) | capitalize"
          :value="item"
          :multiple="true"
        />
      </v-flex>
      <v-btn
        color="primary"
        :disabled="!selected_activities || selected_activities.length == 0"
        block
        @click="submit"
      >
        Save
      </v-btn>
    </v-layout>
  </v-container>
</template>

<script>
  import { MODIFY_PROFILE, SAVE_PROFILE } from 'store/profile/actions'

  import { mapMutations, mapActions, mapState } from 'vuex'

  export default {
    name: 'Activities',
    computed: {
      ...mapState({ activities: (state) => state.profile.activities }),
      selected_activities: {
        set(val) {
          this.MODIFY_PROFILE({ activities: val })
        },
        get() {
          return this.$store.state.profile.profile.activities
        }
      }
    },
    methods: {
      ...mapActions([SAVE_PROFILE]),
      ...mapMutations([MODIFY_PROFILE]),
      submit() {
        this.$loader(true)
        this.SAVE_PROFILE()
          .then(() => {
            this.$loader(false)
            this.$snackbar(
              this.$t('message.success_saving', [this.$t('profile')])
            )
          })
          .catch(() => {
            this.$loader(false)
            this.$snackbar(
              this.$t('message.error_saving', [this.$t('profile')])
            )
          })
      }
    }
  }
</script>
