<template>
  <div class="tabs">
    <div class="headers">
      <Signup />
      <Login />
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import Login from "@/components/Login.vue";
import Signup from "@/components/Signup.vue";

@Component({
  components: {
    Login,
    Signup
  }
})
export default class Tabs extends Vue {
}
</script>