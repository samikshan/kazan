import store from "@/store";
import { VuexModule, Module, Mutation, Action, getModule, MutationAction } from 'vuex-module-decorators'
import { User, UserCreate } from '../models';
import { createUser, hedgehog, createThreadsClient, createBucketsClient, createBucket } from '../api';
import { Buckets, Client } from '@textile/hub';
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
  threadsClient?: Client
  buckets?: Buckets
  bucketKey?: string

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

  get hedgehogIdentity() {
    const wallet = hedgehog.getWallet();
    const identity: HedgehogIdentity = new HedgehogIdentity(wallet.getPrivateKey())
    return HedgehogIdentity
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
  setBuckets(buckets: Buckets) { this.buckets = buckets }

  @Mutation
  setBucketKey(bucketKey: string) { this.bucketKey = bucketKey }

  @Action({commit: 'setBucketKey'})
  async setupUser() {
    try {
      const wallet = hedgehog.getWallet()
      const privKeyBuf = wallet.getPrivateKey()
      // const pubKeyBytes = wallet.getPublicKey()
      const key = await keys.supportedKeys.secp256k1.unmarshalSecp256k1PrivateKey(privKeyBuf)
      // const privKey = await keys.supportedKeys.ed25519.unmarshalEd25519PrivateKey(privKeyBytes)
      console.log(key.public)
      const identity: HedgehogIdentity = new HedgehogIdentity(key)
      console.log(identity);

      const client = await createThreadsClient(identity);
      console.log(client);
      const threadsList = await client.listThreads()
      console.log(threadsList)

      const bucketsClient = await createBucketsClient(identity);
      console.log(bucketsClient);
      const bucketKey = await createBucket(bucketsClient, "kazan-test-bucket")
      const bucketsList = await bucketsClient.list()
      console.log(bucketsList);

      return bucketKey
    } catch (e) {
      console.error(e);
    }

    // const { buckets, bucketKey } = await getBucketKey(identity);

    // console.log(bucketKey);

    // return bucketKey;

    // this.context.commit('setBucketKey', bucketKey)
    // this.context.commit('setBuckets', buckets)
  }
}

export default getModule(UsersModule);