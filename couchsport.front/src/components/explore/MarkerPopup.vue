<template>
  <div class="marker-popup">
    <v-card flat elevation-0 flex d-flex column>
      <img
        v-if="image"
        :src="image.url"
        :alt="image.alt"
        :max-height="image.height"
        :max-width="image.width"
        :width="image.width"
        aspect-ratio="1"
        class="popup-image grey lighten-2"
      />
      <v-card-title class="pa-0">
        <div>
          <h3 class="headline mb-0">
            {{ name }}
          </h3>
        </div>
      </v-card-title>
      <v-divider />
      <v-card-text class="grow pa-0 pt-2 pb-2">
        <div class="pa-0 mb-2">
          {{ desc }}
        </div>

        <div v-if="activities.length > 0">
          <v-chip
            v-for="(a, i) in activities"
            :key="i"
            color="primary"
            text-color="white"
            small
          >
            {{ a.name }}
          </v-chip>
        </div>
      </v-card-text>
      <v-divider />
      <v-card-actions>
        <v-btn class="mr-2" color="primary" small text :href="url">
          {{ $t('p.explore.view_details') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
  export default {
    name: 'MarkerPopup',
    props: {
      name: { type: String, default: '' },
      desc: { type: String, default: '' },
      image: { type: Object, default: () => {} },
      id: { type: Number, default: 0 },
      url: { type: String, default: '' },
      activities: { type: Array, default: () => [] }
    }
  }
</script>

<style lang="scss">
  .leaflet-popup {
    min-width: 293px;
    .marker-popup {
      max-width: 255px;
      .popup-image {
        border-radius: 7px !important;
      }
    }
  }
</style>
