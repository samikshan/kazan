<template>
  <v-card
    class="d-flex align-center justify-center pa-4 mx-auto"
    flat
    outlined>
      <v-card-title>
        <h1 class="display-1">{{ title }}</h1>
      </v-card-title>
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";

@Component
export default class PageHeader extends Vue {
  @Prop() title!: string
}
</script>