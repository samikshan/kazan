import store from "@/store";
import {
  VuexModule,
  Module,
  Mutation,
  Action,
  getModule
} from "vuex-module-decorators";
import {
  User,
  UserCreate,
  Track,
  RecordedTrack,
  UserTrackIndex
} from "@/store/models";
import {
  createUser,
  hedgehog,
  createBucketsClient,
  createBucket,
  getTrackIndex,
  uploadTrackToBucket,
  getTracks,
  storeIndex
} from "@/store/api";
import { HedgehogIdentity } from "@/store/identity";
import { keys } from "libp2p-crypto";
import { Buckets } from "@textile/hub";

@Module({
  namespaced: true,
  name: "users",
  dynamic: true,
  store: store
})
class UsersModule extends VuexModule {
  user: User | null = null;
  buckets: Buckets | null = null;
  bucketKey = "";
  trackIndex: UserTrackIndex | null = null;
  tracks: Array<Track> = [];

  get username() {
    return (this.user && this.user.username) || null;
  }

  get isLoggedIn() {
    const wallet = hedgehog.getWallet();
    console.log(wallet);
    return hedgehog.isLoggedIn();
  }

  get userBucketKey() {
    console.log(this.bucketKey);
    return this.bucketKey;
  }

  get getTracks() {
    return this.tracks;
  }

  @Mutation
  setUser(user: User) {
    this.user = user;
  }

  @Action({ commit: "setUser" })
  async signup(userCreateReq: UserCreate) {
    await createUser(userCreateReq);
    const wallet = hedgehog.getWallet();
    const user: User = {
      username: userCreateReq.username,
      walletAddr: wallet
    };

    return user;
  }

  @Mutation
  setBucketKey(bucketKey: string) {
    this.bucketKey = bucketKey;
  }

  @Mutation
  setBuckets(buckets: Buckets) {
    this.buckets = buckets;
  }

  @Mutation
  setTrackIndex(trackIndex: UserTrackIndex) {
    this.trackIndex = trackIndex;
  }

  @Mutation
  setTracks(tracks: Track[]) {
    this.tracks = tracks;
  }

  @Action
  async setupUser() {
    try {
      const wallet = hedgehog.getWallet();
      const privKeyBuf = wallet.getPrivateKey();
      const key = await keys.supportedKeys.secp256k1.unmarshalSecp256k1PrivateKey(
        privKeyBuf
      );
      const identity: HedgehogIdentity = new HedgehogIdentity(key);

      const bucketsClient = await createBucketsClient(identity);
      const bucketKey = await createBucket(bucketsClient, "kazan-test-bucket");
      const trackIndex = await getTrackIndex(
        bucketsClient,
        bucketKey,
        identity
      );
      console.log(trackIndex);
      const tracks = await getTracks(bucketsClient, bucketKey, trackIndex);
      console.log(tracks);
      this.context.commit("setBucketKey", bucketKey);
      this.context.commit("setBuckets", bucketsClient);
      this.context.commit("setTrackIndex", trackIndex);
      this.context.commit("setTracks", tracks);
    } catch (e) {
      console.error(e);
    }
  }

  @Mutation
  updateTracks(track: Track) {
    this.tracks.push(track);
  }

  @Mutation
  async updateTrackIndex(path: string) {
    if (!this.trackIndex || !this.buckets) {
      throw new Error("track index or buckets is undefined");
    }
    this.trackIndex.paths.push(path);
    await storeIndex(this.trackIndex, this.buckets, this.bucketKey);
  }

  @Action
  async addNewTrack(recordedTrack: RecordedTrack) {
    console.log("about to push new track to bucket");
    if (!this.buckets) {
      throw new Error("bucket client not yet initialised");
    }

    if (!this.trackIndex) {
      throw new Error("track index not yet created for user");
    }
    try {
      const resp = await uploadTrackToBucket(
        recordedTrack,
        this.bucketKey,
        this.buckets
      );
      const track: Track = {
        src: `https://${resp.cid}.ipfs.hub.textile.io`,
        name: resp.name,
        cid: resp.cid
      };
      this.context.commit("updateTrackIndex", resp.metapath);
      this.context.commit("updateTracks", track);
    } catch (err) {
      console.error(err);
    }
  }
}

export default getModule(UsersModule);
