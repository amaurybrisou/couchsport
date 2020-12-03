<template>
  <v-snackbar
    v-model="show"
    :bottom="y === 'bottom'"
    :left="x === 'left'"
    :multi-line="mode === 'multi-line'"
    :right="x === 'right'"
    :top="y === 'top'"
    :vertical="mode === 'vertical'"
    :timeout="3000"
  >
    {{ local_text }}
    <v-btn color="warning" text @click="show = false">
      {{ $t('close') }}
    </v-btn>
  </v-snackbar>
</template>

<script>
  export default {
    name: 'AppSnackBar',
    props: {
      state: { type: Boolean, default: false },
      y: { type: String, default: 'top' },
      x: { type: String, default: null },
      mode: { type: String, default: '' },
      text: {
        type: String,
        default: 'Your profile has been successfully saved'
      }
    },
    data() {
      return {
        show: false,
        color: '',
        local_text: this.text || '',
        timeout: 0
      }
    },
    created() {
      this.$on('snackbar:show', function (state) {
        this.local_text = state.text
        this.color = state.color
        this.timeout = state.timeout
        this.show = true
      })
    }
  }
</script>
