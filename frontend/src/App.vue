<template>
  <v-app>
    <Nav />
    <v-main>
      <!-- Provides the application the proper gutter -->
      <v-container fluid>
        <PageHeader v-bind:title="pageTitle"></PageHeader>
        <!-- If using vue-router -->
        <router-view></router-view>
      </v-container>
    </v-main>
    <!-- <router-view></router-view> -->
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Nav from "@/components/TheNavigation.vue";
import PageHeader from "@/components/ThePageHeader.vue";
@Component({
  components: {
    Nav,
    PageHeader
  }
})
export default class App extends Vue {
  get pageTitle() {
    return this.$route.meta.title
  }
}
</script>