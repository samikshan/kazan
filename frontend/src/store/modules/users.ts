import store from "@/store";
import { VuexModule, Module, Mutation, Action, getModule, MutationAction } from 'vuex-module-decorators'
import { User, Profile } from '../models';
import { getStoredIdentity, createIdentity } from '../api';

@Module({
  namespaced: true,
  name: 'users',
  dynamic: true,
  store: store
})
class UsersModule extends VuexModule {
  user: User | null = null
  profile: Profile | null = null

  get userIdentity() {
    if (this.user) {
      return this.user.identity;
    }

    return getStoredIdentity();
  }

  @Mutation
  setUser(user: User) { this.user = user }

  @Action({commit: 'setUser'})
  async login() {
    const identity = await createIdentity();
    const user: User = {
      identity: identity
    }

    return user
  }
}

export default getModule(UsersModule);