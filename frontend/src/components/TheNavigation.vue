<template>
  <v-navigation-drawer app
    absolute
  >
    <v-img src="../assets/kazan-logo-1.png">
      <!-- <v-row align="end" class="lightbox white--text pa-2 fill-height">
        <v-col>
          <div class="subheading">Jonathan Lee</div>
          <div class="body-1">heyfromjonathan@gmail.com</div>
        </v-col>
      </v-row> -->
    </v-img>

    <v-container>
      <v-list
        nav
        class="py-0"
      >
        <v-list-item>
          <v-list-item-content v-if="!displayName">
            <v-btn
              text
              small
              @click="handleCreateAccount()"
            >
              Create an account
            </v-btn>
          </v-list-item-content>
          <v-list-item-content v-else>
            <v-btn
              text
              small
              @click="handleOpenProfile(userID)"
            >
              {{ displayName }}
            </v-btn>
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
      
        <v-spacer></v-spacer>
        
        <v-btn
          class="ma-2"
          @click="$router.push('record')"
        >   
          Record Tracks
          <v-icon right>mdi-cloud-upload</v-icon>
        </v-btn>
      </v-list>
    </v-container>
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
        await users.getLoggedInUser();
        await users.setupUserBuckets();
        await users.loadUserFeed();
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

  get displayName() {
    return users.displayName;
  }

  get userID() {
    return users.userID;
  }

  handleCreateAccount() {
    this.$router.push({ name: "signup" });
  }

  handleOpenProfile(id: number) {
    this.$router.push({ name: "profile", params: { id: id.toString() }});
  }
}
</script>