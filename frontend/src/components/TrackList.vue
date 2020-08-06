<template>
  <v-container fluid>
    <v-row>
      <v-col
        v-for="(track, index) in tracks"
        :key="index"
        cols="12"
      >
        <v-card color="blue-grey">
          <Track :title="track.Title" :cid="track.CID" />
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import Track from "@/components/Track.vue";
import { TrackMetadata } from "@/store/models";
import users from "@/store/modules/users";

@Component({
  components: {
    Track,
  }
})
export default class TrackList extends Vue {
  @Prop() tracks!: Array<TrackMetadata>;
}
</script>