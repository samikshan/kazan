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
  UserUpdate,
  TrackMetadata,
  RecordedTrack,
  UserTrackIndex,
  StoreTrackMetadataResp,
  StoreTrackMetadata,
} from "@/store/models";
import {
  createUser,
  hedgehog,
  createBucketsClient,
  createBucket,
  uploadTrackToBucket,
  storeTrackFn,
  updateUserFn,
  getUserFn
} from "@/store/api";
import { HedgehogIdentity } from "@/store/hedgehogIdentity";
import { TextileIdentity } from "@/store/textileidentity";
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
  tracks: Array<TrackMetadata> = [];

  get username() {
    return (this.user && this.user.username) || null;
  }

  get walletAddr() {
    return (this.user && this.user.walletAddr) || null;
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
    try {
      await createUser(userCreateReq)
    
      const wallet = hedgehog.getWallet();
      const privKeyBuf = wallet.getPrivateKey();
      const identity: HedgehogIdentity = new HedgehogIdentity(privKeyBuf);
      const respData: any = await getUserFn(identity);
      console.log(respData);

      const user: User = {
        id: respData.id,
        username: respData.username,
        walletAddr: respData.walletAddress,
        instruments: []
      };

      return user;
    } catch (err) {
      console.error(err);
    }
  }

  @Action({ commit: "setUser" })
  async update(userUpdateReq: UserUpdate) {
    try {
      if (!this.user) {
        throw new Error("no logged in user")
      }
      const wallet = hedgehog.getWallet();
      const privKeyBuf = wallet.getPrivateKey();
      const identity: HedgehogIdentity = new HedgehogIdentity(privKeyBuf);

      const respData: User = await updateUserFn(this.user.id, userUpdateReq, identity);
      
      const user: User = {
        id: respData.id,
        username: respData.username,
        walletAddr: respData.walletAddr,
        instruments: []
      };
      return user;
    } catch(err) {
      console.error(err);
    }

    return this.user;
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
  setTracks(tracks: TrackMetadata[]) {
    this.tracks = tracks;
  }

  @Action
  async setupUser() {
    try {
      const wallet = hedgehog.getWallet();

      console.log(wallet.getPublicKeyString());

      const privKeyBuf = wallet.getPrivateKey();
      const key = await keys.supportedKeys.secp256k1.unmarshalSecp256k1PrivateKey(
        privKeyBuf
      );

      const identity: TextileIdentity = new TextileIdentity(key);

      console.log(identity.public.toString());

      const bucketsClient = await createBucketsClient(identity);
      const bucketKey = await createBucket(bucketsClient, "kazan-newest-bucket");

      this.context.commit("setBucketKey", bucketKey);
      this.context.commit("setBuckets", bucketsClient);
    } catch (e) {
      console.error(e);
    }
  }

  @Mutation
  updateTracks(track: TrackMetadata) {
    this.tracks.push(track);
  }

  @Action
  async addNewTrack(recordedTrack: RecordedTrack) {
    console.log("about to push new track to bucket");
    if (!this.buckets) {
      throw new Error("bucket client not yet initialised");
    }
    try {
      const resp = await uploadTrackToBucket(
        recordedTrack,
        this.bucketKey,
        this.buckets
      );

      const trackMetadata: StoreTrackMetadata = {
        cid: resp.cid,
        title: resp.name,
        parentTrackID: recordedTrack.parentTrackID,
        components: []
      }

      const respData: StoreTrackMetadataResp = await storeTrackFn(trackMetadata);

      const track: TrackMetadata = {
        cid: respData.cid,
        title: respData.title,
        composerID: respData.composerID,
        composer: respData.composer,
        parentTrackID: respData.parentTrackID,
        parentTrack: respData.parentTrack,
        forks: respData.forks,
        components: respData.components
      };
      // this.context.commit("updateTrackIndex", resp.metapath);
      this.context.commit("updateTracks", track);
    } catch (err) {
      console.error(err);
    }
  }
}

export default getModule(UsersModule);
