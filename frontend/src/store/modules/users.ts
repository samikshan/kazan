import store from "@/store";
import { VuexModule, Module, Mutation, Action, getModule, MutationAction } from 'vuex-module-decorators'
import { User, UserCreate } from '../models';
import { createUser, hedgehog, createThreadsClient, createBucketsClient, createBucket } from '../api';
import { HedgehogIdentity } from '../identity';
import { keys } from 'libp2p-crypto';

@Module({
  namespaced: true,
  name: 'users',
  dynamic: true,
  store: store
})
class UsersModule extends VuexModule {
  user: User | null = null;
  bucketKey = ""

  get username() {
    return this.user && this.user.username || null;
  }

  get isLoggedIn() {
    const wallet = hedgehog.getWallet();
    console.log(wallet);
    return hedgehog.isLoggedIn()
  }

  get userBucketKey() {
    console.log(this.bucketKey)
    return this.bucketKey
  }

  @Mutation
  setUser(user: User) { this.user = user }

  @Action({commit: 'setUser'})
  async signup(userCreateReq: UserCreate) {
    await createUser(userCreateReq);
    const wallet = hedgehog.getWallet();
    const user: User = {
      username: userCreateReq.username,
      walletAddr: wallet
    }

    return user
  }

  @Mutation
  setBucketKey(bucketKey: string) { this.bucketKey = bucketKey }

  @Action({commit: 'setBucketKey'})
  async setupUser() {
    try {
      const wallet = hedgehog.getWallet()
      const privKeyBuf = wallet.getPrivateKey()
      const key = await keys.supportedKeys.secp256k1.unmarshalSecp256k1PrivateKey(privKeyBuf)
      const identity: HedgehogIdentity = new HedgehogIdentity(key);

      const bucketsClient = await createBucketsClient(identity);
      const bucketKey = await createBucket(bucketsClient, "kazan-test-bucket")

      return bucketKey
    } catch (e) {
      console.error(e);
    }
  }
}

export default getModule(UsersModule);