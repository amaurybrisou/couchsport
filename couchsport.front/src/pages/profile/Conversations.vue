<template>
  <v-container fluild grid-list-xs>
    <v-flex
      v-if="!conversations || conversations.length === 0"
      xs12
      text-sm-center
      text-xs-center
    >
      <v-alert color="info">
        {{ $t('message.empty', ['conversations']) | capitalize }}
      </v-alert>
    </v-flex>
    <v-layout row wrap>
      <v-flex v-if="conversations">
        <v-list>
          <v-list-group
            v-for="(c, idx) in conversations"
            :key="`conversation-${c.id}`"
            :class="unread[idx] ? unread_class : ''"
            no-action
            prepend-icon="mdi-message"
            @click="$messenger.setMessagesRead(idx)"
          >
            <v-divider :key="c.id" />
            <v-list-item slot="activator">
              <v-list-item-action>
                <v-icon v-if="unread[idx]" color="warning"> mdi-star </v-icon>
              </v-list-item-action>
              <v-chip
                v-if="c.from_id == connected_profile_id && c.to"
                small
                color
                class="subheading"
              >
                {{
                  c.to.username || c.to.firstname || c.to.lastname || c.to.email
                }}
              </v-chip>
              <v-chip
                v-if="c.from_id != connected_profile_id && c.from"
                small
                color
                class="subheading"
              >
                {{
                  c.from.username ||
                  c.from.firstname ||
                  c.from.lastname ||
                  c.from.email
                }}
              </v-chip>
              <v-list-item v-if="c.messages" class="text--primary">
                {{ $t('p.conversations.last_message') | capitalize }} :
                {{
                  c.messages[c.messages.length - 1].date
                    | formatDate('MM/DD/YYYY')
                }}
                {{ $t('at') }}
                {{
                  c.messages[c.messages.length - 1].date | formatDate('HH:mm')
                }}
              </v-list-item>

              <!-- <v-list-item>{{ c.to.username || c.to.firstname || c.to.lastname }}</v-list-item> -->
              <v-list-item-action>
                <v-layout row>
                  <v-flex>
                    <v-btn color="primary" text @click="deleteConversation(c)">
                      <v-icon>mdi-delete</v-icon>
                    </v-btn>
                  </v-flex>
                </v-layout>
              </v-list-item-action>
            </v-list-item>

            <v-list-item v-for="m in c.messages" :key="`message-${m.id}`">
              <v-list-item-avatar v-if="m.from_id == connected_profile_id">
                <img v-if="c.to.avatar" :src="c.to.avatar" :alt="c.to.avatar" />
              </v-list-item-avatar>
              <v-list-item-avatar v-if="m.from_id != connected_profile_id">
                <img
                  v-if="c.to.avatar"
                  :src="c.from.avatar"
                  :alt="c.from.avatar"
                />
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title
                  v-if="m.from_id != connected_profile_id"
                  class
                >
                  {{ c.from.username || m.email }}:
                </v-list-item-title>
                <v-list-item-title
                  v-if="m.from_id == connected_profile_id"
                  class
                >
                  You:
                </v-list-item-title>
                <v-list-item class="body-1">
                  {{ m.text }}
                </v-list-item>
              </v-list-item-content>

              <v-list-item-action v-if="m.from_id != connected_profile_id">
                <v-layout row>
                  <v-flex>
                    <v-btn
                      color="primary"
                      text
                      @click.prevent="openMessageDialog(c)"
                    >
                      <v-icon>mdi-reply</v-icon>
                    </v-btn>
                  </v-flex>
                </v-layout>
              </v-list-item-action>
            </v-list-item>
          </v-list-group>
        </v-list>
      </v-flex>
    </v-layout>
    <v-layout row justify-center>
      <v-dialog
        v-if="conversations"
        id="contact-dialog"
        v-model="showContactDialog"
        width="500"
      >
        <v-card>
          <v-toolbar color="primary">
            <v-card-title class="title font-weight-regular">
              {{ $t('p.conversations.write_your_reply') }}
            </v-card-title>
          </v-toolbar>
          <v-form v-model="messageFormValid">
            <v-card-text>
              <v-text-field
                v-if="!email"
                v-model="message.email"
                name="email"
                :value="$t('email') | capitalize"
                autocomplete="email"
                :rules="emailRules"
                autofocus
              />
              <v-textarea
                v-model="message.text"
                name="Message"
                :placeholder="$t('_message') | capitalize"
                :rules="textRules"
                row="1"
                maxlength="128"
                hide-details
                no-resize
                autofocus
              />
            </v-card-text>
            <v-card-actions>
              <v-spacer />
              <v-btn
                color="primary"
                text
                @click.prevent.native="showContactDialog = false"
              >
                {{ $t('cancel') }}
              </v-btn>
              <v-btn
                color="primary"
                text
                :disabled="!messageFormValid"
                @click.native="reply"
              >
                {{ $t('send') }}
              </v-btn>
            </v-card-actions>
          </v-form>
        </v-card>
      </v-dialog>
    </v-layout>
  </v-container>
</template>

<script>
  import { mapState, mapActions } from 'vuex'
  import {
    GET_CONVERSATIONS,
    REMOVE_CONVERSATION
  } from 'store/conversations/actions'

  const NAMESPACE = 'conversations/'
  export default {
    name: 'Conversations',
    data() {
      return {
        unread_class: 'unread_conversation',

        focusedConversation: null,
        showContactDialog: false,
        messageFormValid: false,
        emailRules: [
          (v) => !!v || this.$t('message.required', ['', this.$t('_message')]),
          (v) =>
            /.+@.+/.test(v) || this.$t('message.invalid', [this.$t('email')])
        ],

        textRules: [
          (v) => !!v || this.$t('message.required', ['', this.$t('_message')]),
          (v) =>
            (v && v.length >= 20) ||
            this.$t('message.length_above', { len: 20 })
        ]
      }
    },
    computed: {
      unread() {
        return this.$store.state.profile.conversations.conversations.map(
          (c) => {
            return c.unread
          }
        )
      },
      conversations: {
        get() {
          return this.$store.state.profile.conversations.conversations
        }
      },
      ...mapState({
        email: (state) => state.auth.email,
        connected_profile_id: (state) => state.profile.profile.id
      }),
      message() {
        return { from_id: null, to_id: null, email: this.email, Text: '' }
      }
    },
    mounted() {
      this.GET_CONVERSATIONS()
    },
    methods: {
      ...mapActions(NAMESPACE, [GET_CONVERSATIONS, REMOVE_CONVERSATION]),
      openMessageDialog: function (c) {
        this.showContactDialog = true
        this.message.to_id =
          c.from_id === this.connected_profile_id ? c.to_id : c.from_id
        this.focusedConversation = c
      },
      reply: function () {
        this.$messenger
          .sendMessage(this.message)
          .then(() => {
            this.$snackbar('Your messages has been sent')
          })
          .catch(() => {
            tihs.$snackbar('An error occured while sending your message')
          })
      },
      async deleteConversation(c) {
        this.focusedConversation = c

        if (c.id != null) {
          this.REMOVE_CONVERSATION(c.id)
            .then(() => {
              this.$snackbar('this conversation has been successfully deleted')
            })
            .catch(() => {
              this.$snackbar('there was an error deleting this conversation')
            })
        }
      }
    }
  }
</script>

<style lang="scss">
  .unread_conversation {
    background: rgba(#607d8b, 0.3);
  }
</style>
