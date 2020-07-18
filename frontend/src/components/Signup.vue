<template>
  <div class="form">
    <div class="fields">
      <input v-model="username" placeholder="Username">
      <input v-model="password" placeholder="Password" type="password">
      <div>
        <input v-model="passwordConfirmation" placeholder="Confirm Password" type="password">
        <!-- <p className="error">{{ errorMessage }}</p> -->
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
    console.log("Signup request!");
    console.log(this.username, this.password);
    if (this.password !== this.passwordConfirmation) {
      console.log("passwords dont match");
      this.errorMessage = "The passwords you entered don't match"
    } else if (!this.username || !this.password || !this.passwordConfirmation) {
      console.log("incomplete fields");
      console.log(this.username, this.password);
      this.errorMessage = "Please enter the required fields"
    } else {
      this.errorMessage = ""
      try {
        const user = await users.signup({
          username: this.username,
          password: this.password
        })
        console.log(user);
      } catch(e) {
        console.error(e);
        this.errorMessage = "Account already exists. Try logging in instead"
      }
    }
  }
}
</script>
