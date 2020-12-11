<template>
  <div>
    <v-layout v-if="pages">
      <v-flex md8 offset-md2>
        <v-list>
          <template v-for="(p, idx) in pages">
            <v-divider v-if="idx < pages.length - 1" :key="idx"></v-divider>
            <v-list-item :key="`preview-image-${p.id}`" class="page-line">
              <v-list-item-avatar>
                <img
                  v-if="p.images && p.images.length > 0"
                  :src="p.images[0].url"
                  :alt="p.images[0].alt"
                />
              </v-list-item-avatar>

              <v-list-item-title class="ml-3">
                {{ p.description | shorten(30) }}
              </v-list-item-title>

              <v-list-item-action class="page-list-item">
                <v-checkbox
                  :input-value="p.public"
                  :color="$vuetify.breakpoint.xsOnly ? '' : 'primary'"
                  :true-value="true"
                  :false-value="false"
                  :label="
                    ($vuetify.breakpoint.xsOnly
                      ? ''
                      : p.public
                      ? $t('public', [])
                      : $t('private', [])) | capitalize
                  "
                  @change="publishPage(p.id, $event)"
                />

                <page-edition-dialog
                  :state="'edit'"
                  :all-activities="allActivities"
                  @page-saved="onPageSaved"
                >
                  <template v-slot:open-btn="{ on }">
                    <v-btn
                      small
                      :color="$vuetify.breakpoint.xsOnly ? '' : 'primary'"
                      :to="{
                        name: 'page-details',
                        params: { page_name: p.name }
                      }"
                      v-on="on"
                    >
                      <v-icon>mdi-eye</v-icon>
                    </v-btn>
                    <v-btn
                      small
                      :color="$vuetify.breakpoint.xsOnly ? '' : 'primary'"
                      v-on="on"
                      @click="editPage(p.id)"
                    >
                      <v-icon>mdi-pencil</v-icon>
                    </v-btn>
                    <v-btn
                      small
                      :color="$vuetify.breakpoint.xsOnly ? '' : 'primary'"
                      @click.stop="deletePage(p.id)"
                    >
                      <v-icon>mdi-delete</v-icon>
                    </v-btn>
                  </template>

                  <span slot="submitText">
                    {{ $t('save') }} {{ $t('modifications') }}
                  </span>
                  <span slot="pageTitle">
                    {{ $t('edit') }} {{ $t('page') }} : {{ p.title }}
                  </span>
                </page-edition-dialog>
              </v-list-item-action>
            </v-list-item>
          </template>
        </v-list>
      </v-flex>
    </v-layout>
    <v-layout>
      <v-flex d-flex>
        <page-edition-dialog
          :state="'new'"
          :all-activities="allActivities"
          @page-saved="onPageSaved"
        >
          <template v-slot:open-btn="{ on }">
            <v-btn block color="success" text v-on="on">
              {{ $t('new') }} {{ $t('page') }}
            </v-btn>
          </template>
          <span slot="pageTitle"> {{ $t('new') }} {{ $t('page') }} </span>
        </page-edition-dialog>
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
  import PageEditionDialog from 'components/profile/PageEditionDialog'

  import {
    GET_PAGES,
    EDIT_PAGE,
    PUBLISH_PAGE,
    DELETE_PAGE,
    NEW_PAGE
  } from 'store/pages/actions'
  import { mapMutations, mapActions, mapState } from 'vuex'

  const NAMESPACE = 'pages/'

  export default {
    name: 'Pages',
    components: { PageEditionDialog },
    computed: {
      ...mapState({
        allActivities: (state) => state.profile.activities,
        pages: (state) => state.profile.pages.pages
      })
    },
    created() {
      this.GET_PAGES()
    },
    methods: {
      ...mapActions(NAMESPACE, [
        GET_PAGES,
        NEW_PAGE,
        PUBLISH_PAGE,
        DELETE_PAGE
      ]),
      ...mapMutations(NAMESPACE, [EDIT_PAGE]),
      onPageSaved(state) {
        if (state) {
          return this.$snackbar(
            this.$t('message.success_saving', [this.$t('page')])
          )
        }

        this.$snackbar({
          text: this.$t('message.error_saving', [this.$t('page')]),
          color: 'error'
        })
      },
      editPage(id) {
        this.EDIT_PAGE(id)
      },
      deletePage(id) {
        if (id != null) {
          var that = this
          this.DELETE_PAGE({ id: id })
            .then(function () {
              this.$snackbar(
                that.$t('message.success_deleting', [that.$t('page')])
              )
            })
            .catch(() => {
              this.$snackbar(
                that.$t('message.error_deleting', [that.$t('page')])
              )
            })
        }
      },
      publishPage(id, state) {
        if (id != null && (state === false || state === true)) {
          this.PUBLISH_PAGE({ id: id, public: state })
            .then(() => {
              this.$snackbar(
                state
                  ? this.$t('message.state', [
                      this.$t('page'),
                      this.$t('public')
                    ])
                  : this.$t('message.state', [
                      this.$t('page'),
                      this.$t('private', ['e'])
                    ])
              )
            })
            .catch(() => {
              this.$snackbar(
                this.$t('message.error_updating', [this.$t('page')])
              )
            })
        }
      }
    }
  }
</script>

<style lang="scss">
  .page-map {
    height: 350px;
  }

  .page-line:hover {
    background: rgba(#607d8b, 0.12);
  }

  .page-list-item {
    display: contents;
    > * {
      margin-left: 0.5vw;
    }
  }
</style>
