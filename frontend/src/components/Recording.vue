<template>
  <v-card>
    <!-- <div class="d-flex flex-no-wrap justify-space-between"> -->
      <!-- <div> -->
      <v-card-title
        class="headline"
        v-text="title"
      ></v-card-title>
      <audio controls v-bind:src="localURL"></audio>
      <!-- </div> -->
    <!-- </div> -->
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";

@Component
export default class Recording extends Vue {
  @Prop() localURL!: string
  @Prop() title!: string
}
</script>