<template>
  <v-app>
    <Signup v-if="!isLoggedIn" />
    <div v-else>
      <p v-if="!bucketKey">Buckets not setup</p>
      <div v-else>
        <p>Buckets set up! Can upload stuff to Textile now</p>
        <Recorder />
        <Profile />
      </div>
    </div>
  </v-app>
</template>

<script lang="ts">
// @ is an alias to /src
import Signup from "@/components/Signup.vue";
import Recorder from "@/components/Recorder.vue";
import Profile from "@/components/Profile.vue";
import users from "@/store/modules/users";
import { Component, Vue } from "vue-property-decorator";

@Component({
  components: {
    Signup,
    Recorder,
    Profile
  }
})
export default class Home extends Vue {
  async created() {
    if (users.isLoggedIn && !users.userBucketKey) {
      try {
        await users.setupUser();
      } catch (e) {
        console.error(e);
      }
    }
  }

  get isLoggedIn() {
    return users.isLoggedIn;
  }

  get bucketKey() {
    return users.userBucketKey;
  }
}
</script>
