<template>
  <v-dialog v-model="show" hide-overlay persistent width="300">
    <v-card color="primary" dark>
      <v-card-text>
        {{ $t(message) }}
        <v-progress-linear indeterminate color="white" class="mb-0" />
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
  export default {
    name: 'AppLoader',
    props: {
      message: {
        type: String,
        default: () => {
          return 'message.stand_by'
        },
        required: false
      }
    },
    data() {
      return {
        show: false
      }
    },
    created() {
      this.$on('loader:toggle', function (state) {
        this.show = state.status
        this.message = state.message
      })
    }
  }
</script>
