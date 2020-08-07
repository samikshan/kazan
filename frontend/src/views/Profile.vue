<template>
  <v-app>
    <v-card
      class="pa-2 ma-2 d-flex align-left"
      dark
      color="grey darken-1"
    >
      <v-card-title>
        {{ profile.displayName }}
      </v-card-title>
    </v-card>
    <v-card
      class="pa-2 ma-2 d-flex align-left"
      dark
      color="grey darken-1"
    >
      <v-card-title>
        Instruments
      </v-card-title>
      <v-row
        v-for="(instrument, i) in profile.instruments"
        :key="i"
      >
        <v-card-text>
          instrument.name
        </v-card-text>   
      </v-row>
    </v-card>
  </v-app>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import profiles from "@/store/modules/profiles";

@Component
export default class Profile extends Vue {
  get profile() {
    return profiles.profile;
  }
  
  async created() {
    try {
      const id: number = parseInt(this.$route.params.id)
      console.log(id);
      await profiles.fetchProfile(id);
    } catch (err) {
      console.error(err);
    }
  }
}
</script>