<template>
  <div class="form">
    <div class="fields">
      <input class="[ errorMessage ? error : null ]" v-model="username" placeholder="Username">
      <input class="[ errorMessage ? error : null ]" v-model="password" placeholder="Password" type="password">
      <div>
        <input class="[ errorMessage ? error : null ]" v-model="passwordConfirmation" placeholder="Confirm Password" type="password">
        <p className="error">{{ errorMessage }}</p>
      </div>
    </div>
    <div className="buttons">
      <button @click="handleSignup()">Create My Account</button>
    </div>
  </div>
</template>

<script lang="ts">

import { Component, Vue } from "vue-property-decorator";
import users from "@/store/modules/users";

@Component
export default class Signup extends Vue {
  username = ""
  password = ""
  passwordConfirmation = ""
  errorMessage = ""

  async handleSignup() {
    console.log(this.username, this.password);
    if (this.password !== this.passwordConfirmation) {
      this.errorMessage = "The passwords you entered don't match"
    } else if (!this.username || !this.password || !this.passwordConfirmation) {
      this.errorMessage = "Please enter the required fields"
    } else {
      this.errorMessage = ""
      try {
        await users.signup({
          username: this.username,
          password: this.password
        })
      } catch(e) {
        console.error(e);
        this.errorMessage = "Account already exists. Try logging in instead"
      }
    }
  }
}
</script>
