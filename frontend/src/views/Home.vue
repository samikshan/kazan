<template>
  <div class="home">
    <!-- <CreateProfile v-if="!username"/> -->
    <AuthTabs v-if="!username" />
  </div>
</template>

<script lang="ts">
// @ is an alias to /src
import AuthTabs from "@/components/AuthTabs.vue";
import users from "@/store/modules/users";
import { Component, Vue } from "vue-property-decorator";

@Component({
  components: {
    AuthTabs
  }
})
export default class Home extends Vue {
  get username() {
    return users.username;
  }
}

</script>
