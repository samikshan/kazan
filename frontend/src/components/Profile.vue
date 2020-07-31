<template>
  <v-app>
    <v-card class="mx-auto mt-5" max-width="600px" min-width="360px">
      <v-card-title>
        <h2>My tracks</h2>
      </v-card-title>

      <v-card v-if="!showMixer">
        <ul class="tracks-list">
          <li v-for="track in tracks" v-bind:key="track.cid">
            {{ track.name }} : {{ track.cid }}
            <audio controls v-bind:src="track.src"></audio>
            <div id="buttons">
              <v-btn v-bind:id="track.cid" @click="jamOnCID = track.cid">
                Jam on this track
              </v-btn>
            </div>
          </li>
        </ul>
      </v-card>
      <v-card v-if="showMixer">
        <Mixer v-bind:parentAudioCID="jamOnCID"/>
      </v-card>
    </v-card>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import users from "@/store/modules/users";
import { Track } from "../store/models";
import Mixer from "@/components/Mixer.vue";

@Component({
  components: {
    Mixer,
  }
})
export default class Profile extends Vue {
  jamOnCID = ""
  get tracks() {
    return users.getTracks;
  }

  get showMixer() {
    return this.jamOnCID.length > 0;
  }
}
</script>
