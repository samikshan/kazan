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
      <v-card-actions :class="`d-flex justify-center`">
          <v-btn
          :loading="signingUp"
          :disabled="signingUp"
          large
          color="blue-grey darken-4"
          class="mx-auto white--text"
          @click="handleSignup()"
        >
          Sign Up
        </v-btn>
      </v-card-actions>
      
      <v-dialog v-model="dialog" persistent width="400px">
        <CreateProfileDialog
          v-bind:address="walletAddr"
          @profile-created="handleProfileCreated()"
        />
      </v-dialog>
    </v-card>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import CreateProfileDialog from "@/components/CreateProfileDialog.vue";
import { User } from "@/store/models";
import users from "@/store/modules/users";

@Component({
  components: {
    CreateProfileDialog
  }
})
export default class Signup extends Vue {
  username = "";
  password = "";
  walletAddr = "";
  confirmPassword = "";
  showPassword = false;
  showConfirmPassword = false;
  errorMessage = "";
  dialog = false;
  signingUp = false;

  async handleSignup() {
    console.log(this.username, this.password);
    if (this.password !== this.confirmPassword) {
      this.errorMessage = "The passwords you entered don't match";
    } else if (!this.username || !this.password || !this.confirmPassword) {
      this.errorMessage = "Please enter the required fields";
    } else {
      this.errorMessage = "";
      try {
        console.log("Signing up...")
        this.signingUp = true;
        const user: any = await users.signup({
          username: this.username,
          password: this.password
        });

        console.log(user);

        users.setupUserBuckets();
        
        this.walletAddr = user.walletAddr;
        this.dialog = true;
        this.signingUp = false;
      } catch (e) {
        console.error(e);
        this.signingUp = false;
        this.errorMessage = "Account already exists. Try logging in instead";
      }
    }
  }

  async handleProfileCreated() {
    try {
      this.dialog = false;
      await users.loadUserFeed();
      this.$router.push("/");
    } catch (err) {
      console.error(err);
    }
  }
}
</script>