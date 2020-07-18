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
    return this.user && this.user.username || null;
  }

  get isLoggedIn() {
    return hedgehog.isLoggedIn()
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