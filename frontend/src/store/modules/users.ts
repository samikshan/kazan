import store from "@/store";
import { VuexModule, Module, Mutation, Action, getModule, MutationAction } from 'vuex-module-decorators'
import { User, Profile } from '../models';
import { getIdentity } from '../api';

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
    return this.user && this.user.identity || null;
  }

  @Mutation
  setUser(user: User) { this.user = user }

  @Action
  async login() {
    const identity = await getIdentity();
    const user: User = {
      identity: identity
    }
    
    return user
  }
}

export default getModule(UsersModule);