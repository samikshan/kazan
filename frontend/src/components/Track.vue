<template>
  <v-card
    class="d-flex align-center justify-space-around mb-6"
    dark
    color="grey darken-1"
  >
    <v-card-title
      class="caption"
      v-text="title"
    ></v-card-title>
    <v-spacer />
    <audio controls v-bind:src="url"></audio>
    <v-spacer />
    <v-btn @click="jamOnTrack()">
      Jam On This One
    </v-btn>
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";

@Component
export default class Track extends Vue {
  @Prop() cid!: string
  @Prop() title!: string

  get url() {
      return `https://${this.cid}.ipfs.hub.textile.io`
  }

  jamOnTrack() {
    console.log("need to jam on this piece of music");
  }
}
</script>