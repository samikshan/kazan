import store from "@/store";
import {
  VuexModule,
  Module,
  Mutation,
  Action,
  getModule
} from "vuex-module-decorators";
import { getProfileFn } from "@/store/api";
import { Profile } from '@/store/models';


@Module({
  namespaced: true,
  name: "profiles",
  dynamic: true,
  store: store
})
class ProfilesModule extends VuexModule {
  profile: Profile = {
    username: "",
    displayName: "",
    instruments: [],
    followingCount: 0,
    followerCount: 0,
    jamCount: 0,
    trackCount: 0
  }

  @Mutation
  setProfile(profile: Profile) {
    this.profile = profile;
  }

  @Action({ commit: "setProfile" })
  async fetchProfile(id: number) {
    const profile: Profile = await getProfileFn(id);
    return profile;
  }
}

export default getModule(ProfilesModule);