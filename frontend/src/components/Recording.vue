<template>
  <v-card
    class="d-flex align-center justify-space-around ma-2 pa-2"
    dark
    color="grey darken-2"
  >
    <audio controls v-bind:src="localURL"></audio>
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