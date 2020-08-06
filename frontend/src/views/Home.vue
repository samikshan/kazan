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
      <v-container fluid v-if="userFeed">
        <v-row dense>
          <v-col
            v-for="(section, index) in userFeed.tracks"
            :key="index"
            cols="12"
          >
            <v-subheader
              :key="section.name"
            >
              {{ section.name }}
            </v-subheader>
            <TrackList v-bind:tracks="section.tracks" />
          </v-col>
        </v-row>
      </v-container>
    </v-container>
  </v-app>
</template>

<script lang="ts">
import users from "@/store/modules/users";
import { Component, Vue, Watch } from "vue-property-decorator";
import TrackList from "@/components/TrackList.vue";
import { TrackMetadata, UserFeed } from '../store/models';

@Component({
  components: {
    TrackList,
  }
})
export default class Home extends Vue {
  searchText? = "";
  recommendedTracks: Array<TrackMetadata> = [];

  get userFeed() {
    return users.getUserFeed;
  }

  handleSearch() {
    if (!this.searchText) {
      console.log("empty search text");
    } else {
      this.$router.push({ name: "search", params: { text: this.searchText }});
    }
  }
}
</script>
