<template>
  <div class="form">
    <div class="fields">
      <input v-model="username" placeholder="Username" />
      <div>
        <input
          class="[ error: errorMessage ? null ]"
          v-model="password"
          placeholder="Password"
          type="password"
        />
        <p className="error">{{ errorMessage }}</p>
      </div>
    </div>
    <div className="buttons">
      <button v-on:click="handleLogin">Log In</button>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";

@Component
export default class Login extends Vue {
  username = "";
  password = "";
  errorMessage = "";

  handleLogin = async (event: Event) => {
    console.log("Login request!");
  };
}
</script>
