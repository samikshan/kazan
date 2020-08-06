<template>
  <v-card
    class="d-flex align-center justify-space-around mb-6"
    dark
    color="grey darken-2"
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
    this.$router.push({ name: "jam", params: { parentCID: this.cid } });
  }
}
</script>