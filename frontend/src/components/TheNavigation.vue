<template>
  <v-navigation-drawer app
    absolute
  >
    <v-img src="../assets/kazan-logo-1.png">
      <v-row align="end" class="lightbox white--text pa-2 fill-height">
        <v-col>
          <div class="subheading">Jonathan Lee</div>
          <div class="body-1">heyfromjonathan@gmail.com</div>
        </v-col>
      </v-row>
    </v-img>

    <v-list
      dense
      nav
      class="py-0"
    >
      <v-list-item :class="miniVariant && 'px-0'">
        <v-list-item-content v-if="!username">
          <router-link :to="{ name: 'signup' }">
            Create an account
          </router-link>
        </v-list-item-content>
        <v-list-item-content v-else>
          <router-link :to="'/@' + username">
            {{username}}
          </router-link>
        </v-list-item-content>
      </v-list-item>

      <v-divider></v-divider>

      <v-list-item
        v-for="item in items"
        :key="item.title"
        link
      >
        <v-list-item-content>
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import users from "@/store/modules/users";

@Component
export default class Nav extends Vue {
  items = [
    { title: 'My Jams' },
    { title: 'Feed' }
  ]

  async created() {
    if (users.isLoggedIn && !users.userBucketKey) {
      try {
        await users.setupUser();
      } catch (e) {
        console.error(e);
      }
    }
  }

  get isLoggedIn() {
    return users.isLoggedIn;
  }

  get bucketKey() {
    return users.userBucketKey;
  }

  get username() {
    return users.username;
  }
}
</script>