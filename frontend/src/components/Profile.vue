<template>
  <div class="record-upload-track">
    <h2>My tracks</h2>

    <ul class="tracks-list">
      <li v-for="track in tracks" v-bind:key="track.cid">
        {{ track.name }} : {{ track.cid }}
        <audio
          controls
          v-bind:src="track.src"
        ></audio>
      </li>
    </ul>
  </div>
</template>

<script lang="ts">

import { Component, Vue } from "vue-property-decorator";
import users from "@/store/modules/users";
import { Track } from '../store/models';

@Component
export default class Profile extends Vue {
  get tracks() {
    return users.getTracks
  }
}

</script>