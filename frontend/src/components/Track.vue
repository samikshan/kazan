<template>
  <v-card>
    <!-- <div class="d-flex flex-no-wrap justify-space-between"> -->
      <!-- <div> -->
      <v-card-title
        class="headline"
        v-text="title"
      ></v-card-title>
      <audio controls v-bind:src="url"></audio>
      <!-- </div> -->
    <!-- </div> -->
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
}
</script>