import store from "@/store";
import { VuexModule, Module, Mutation, Action, getModule, MutationAction } from 'vuex-module-decorators'
import { User, UserCreate } from '../models';
import { getStoredIdentity, createIdentity, createUser, hedgehog } from '../api';

@Module({
  namespaced: true,
  name: 'users',
  dynamic: true,
  store: store
})
class UsersModule extends VuexModule {
  user: User | null = null;

  get username() {
    // if (this.user) {
    //   return this.user.identity;
    // }
    return this.user && this.user.username || null;
  }

  @Mutation
  setUser(user: User) { this.user = user }

  @Action({commit: 'setUser'})
  async signup(userCreateReq: UserCreate) {
    await createUser(userCreateReq);
    const user: User = {
      username: userCreateReq.username,
      walletAddr: hedgehog.getWallet()
    }

    return user
  }
}

export default getModule(UsersModule);