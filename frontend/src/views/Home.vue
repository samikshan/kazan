<template>
  <div class="home">
    <CreateProfile v-if="!userIdentity"/>
  </div>
</template>

<script lang="ts">
// @ is an alias to /src
import CreateProfile from "@/components/CreateProfile.vue";
import users from "@/store/modules/users";
import { Component, Vue } from "vue-property-decorator";

@Component({
  components: {
    CreateProfile
  }
})
export default class Home extends Vue {
  get userIdentity() {
    return users.userIdentity;
  }
}

</script>
