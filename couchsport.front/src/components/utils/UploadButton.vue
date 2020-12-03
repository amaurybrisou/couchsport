<template>
  <div class="upload-button">
    <div @click="onFocus">
      <slot name="appearance">
        <v-text-field
          ref="fileTextField"
          v-model="filename"
          class="mt-0 pt-0"
          prepend-icon="mdi-file"
          single-line
          :label="label"
          :required="required"
          :disabled="disabled"
          :error-messages="errors"
          :error="errors.length > 0"
        />
      </slot>
    </div>
    <input
      ref="fileInput"
      type="file"
      :accept="accept"
      :multiple="false"
      :disabled="disabled"
      @change="onFileChange"
    />
  </div>
</template>

<script>
  export default {
    name: 'UploadButton',
    props: {
      value: {
        type: [Array, String],
        default: () => []
      },
      accept: {
        type: String,
        default: '*'
      },
      label: {
        type: String,
        default: 'Please choose...'
      },
      required: {
        type: Boolean,
        default: false
      },
      disabled: {
        type: Boolean,
        default: false
      },
      multiple: {
        type: Boolean, // not yet possible because of data
        default: false
      },
      errors: {
        type: Array,
        default: () => []
      }
    },
    data() {
      return {
        filename: ''
      }
    },
    watch: {
      value(v) {
        this.filename = v
      }
    },
    mounted() {
      this.filename = this.value
    },

    methods: {
      getFormData(files) {
        const data = new FormData()
        ;[...files].forEach((file) => {
          data.append('file', file, file.name) // currently only one file at a time
        })
        return data
      },
      onFocus() {
        if (!this.disabled) {
          this.$refs.fileInput.click()
        }
      },
      onFileChange($event) {
        const files = $event.target.files || $event.dataTransfer.files
        const form = this.getFormData(files)
        if (files) {
          if (files.length > 0) {
            this.filename = [...files].map((file) => file.name).join(', ')
          } else {
            this.filename = null
          }
        } else {
          this.filename = $event.target.value.split('\\').pop()
        }
        this.$emit('input', this.filename)
        this.$emit('formData', form)
      }
    }
  }
</script>

<style lang="scss" scoped>
  input[type='file'] {
    position: absolute;
    left: -99999px;
  }
</style>
