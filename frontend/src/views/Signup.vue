<template>
  <v-app>
    <v-card class="mx-auto mt-5" max-width="600px" min-width="360px">
      <v-card-title>
        <h1 class="display-1">Signup for Kazan</h1>
      </v-card-title>
      <v-card-text>
        <v-form>
          <v-text-field
            label="Email"
            v-model="username"
            prepend-icon="mdi-account-circle"
          />
          <v-text-field
            :type="showPassword ? 'text' : 'password'"
            label="Password"
            v-model="password"
            prepend-icon="mdi-lock"
            :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
            @click:append="showPassword=!showPassword"
          />
          <v-text-field
            :type="showConfirmPassword ? 'text' : 'password'"
            label="Confirm password"
            v-model="confirmPassword"
            prepend-icon="mdi-lock"
            :append-icon="showConfirmPassword ? 'mdi-eye' : 'mdi-eye-off'"
            @click:append="showConfirmPassword=!showConfirmPassword"
          />
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="handleSignup()">Create My Account</v-btn>
      </v-card-actions>
    </v-card>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import users from "@/store/modules/users";

@Component
export default class Signup extends Vue {
  username = "";
  password = "";
  confirmPassword = "";
  showPassword = false;
  showConfirmPassword = false;
  errorMessage = "";

  async handleSignup() {
    console.log(this.username, this.password);
    if (this.password !== this.confirmPassword) {
      this.errorMessage = "The passwords you entered don't match";
    } else if (!this.username || !this.password || !this.confirmPassword) {
      this.errorMessage = "Please enter the required fields";
    } else {
      this.errorMessage = "";
      try {
        users.signup({
          username: this.username,
          password: this.password
        }).then(() => this.$router.push("/"));
      } catch (e) {
        console.error(e);
        this.errorMessage = "Account already exists. Try logging in instead";
      }
    }
  }
}
</script>