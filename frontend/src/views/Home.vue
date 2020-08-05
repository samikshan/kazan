<template>
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
import users from "@/store/modules/users";
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Home extends Vue {
  searchText? = "";

  handleSearch() {
    if (!this.searchText) {
      console.log("empty search text")
    } else {
      this.$router.push({ name: "search", params: { text: this.searchText }});
    }
  }
}
</script>
