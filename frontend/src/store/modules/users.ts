import store from "@/store";
import { VuexModule, Module, Mutation, Action, getModule, MutationAction } from 'vuex-module-decorators'
import { User, UserCreate, Track, RecordedTrack } from '@/store/models';
import { createUser, hedgehog, createThreadsClient, createBucketsClient, createBucket, pushToBucket } from '@/store/api';
import { HedgehogIdentity } from '@/store/identity';
import { keys } from 'libp2p-crypto';
import { Buckets } from '@textile/hub';

@Module({
  namespaced: true,
  name: 'users',
  dynamic: true,
  store: store
})
class UsersModule extends VuexModule {
  user: User | null = null;
  buckets: Buckets | null = null;
  bucketKey = ""
  tracks: Array<Track> = []

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

  get getTracks() {
    return this.tracks
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

  @Mutation
  setBuckets(buckets: Buckets) { this.buckets = buckets }

  @Action
  async setupUser() {
    try {
      const wallet = hedgehog.getWallet()
      const privKeyBuf = wallet.getPrivateKey()
      const key = await keys.supportedKeys.secp256k1.unmarshalSecp256k1PrivateKey(privKeyBuf)
      const identity: HedgehogIdentity = new HedgehogIdentity(key);

      const bucketsClient = await createBucketsClient(identity);
      const bucketKey = await createBucket(bucketsClient, "kazan-test-bucket")

      this.context.commit('setBucketKey', bucketKey)
      this.context.commit('setBuckets', bucketsClient)
    } catch (e) {
      console.error(e);
    }
  }

  @Mutation
  updateTracks(track: Track) { this.tracks.push(track) }

  @Action({commit: 'updateTracks'})
  async addNewTrack(recordedTrack: RecordedTrack) {
    console.log("about to push new track to bucket")
    if (!this.buckets) {
      throw new Error("bucket client not yet initialised");
    }
    const raw = await pushToBucket(recordedTrack, this.bucketKey, this.buckets);
    const track: Track = {
      cid: raw.path.cid.toString(),
      name: recordedTrack.name,
      path: raw.path.path
    }
    return track
  }
}

export default getModule(UsersModule);