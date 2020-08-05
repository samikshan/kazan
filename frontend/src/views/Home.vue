<template>
    <!-- <Signup v-if="!isLoggedIn" />
    <div v-else>
      <p v-if="!bucketKey">Buckets not setup</p>
      <div v-else>
        <p>Buckets set up! Can upload stuff to Textile now</p>
        <Recorder />
        <Profile />
      </div>
    </div> -->
  <v-app>
    <v-container fluid>
      <v-text-field
        hide-details
        single-line
        clearable
        v-model="searchText"
        label="Search tracks, users, etc..."
        @keyup.enter="handleSearch"
      ></v-text-field>
      <v-divider></v-divider>
      <v-btn
        @click="$router.push('record')">
          Record And Upload Fresh Tracks
          <v-icon>mdi-cloud-upload</v-icon>
      </v-btn>
    </v-container>
  </v-app>
</template>

<script lang="ts">
// @ is an alias to /src
// import Signup from "@/components/Signup.vue";
// import Recorder from "@/components/Recorder.vue";
// import Profile from "@/components/Profile.vue";
import users from "@/store/modules/users";
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Home extends Vue {
  searchText? = "";
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

  handleSearch() {
    if (!this.searchText) {
      console.log("empty search text")
    } else {
      this.$router.push({ name: "search", params: { text: this.searchText }})
    }
  }
}
</script>
