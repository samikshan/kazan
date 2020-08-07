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
  StoreTrackMetadata,
  Instrument,
  TracksByInstrument,
  UserFeed,
} from "@/store/models";
import {
  createUser,
  hedgehog,
  createBucketsClient,
  createBucket,
  uploadTrackToBucket,
  storeTrackFn,
  updateUserFn,
  getUserFn,
  getUserFeedFn
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
  userFeed: UserFeed | null = null; 
  buckets: Buckets | null = null;
  bucketKey = "";
  trackIndex: UserTrackIndex | null = null;
  tracks: Array<TrackMetadata> = [];
  recommendedTracks: Array<TrackMetadata> = [];
  textileIdent: TextileIdentity | null = null;
  hedgehogIdent: HedgehogIdentity | null = null;

  get username() {
    return (this.user && this.user.username) || null;
  }

  get userID() {
    return (this.user && this.user.id) || null;
  }

  get displayName() {
    return (this.user && this.user.displayName) || null;
  }

  get walletAddr() {
    return (this.user && this.user.walletAddr) || null;
  }

  get isUserIdentSet() {
    return !this.hedgehogIdent;
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

  get getUserFeed() {
    return this.userFeed;
  }

  @Mutation
  setUser(user: User) {
    this.user = user;
  }


  @Mutation
  setHedgehogIdent(identity: HedgehogIdentity) {
    this.hedgehogIdent = identity;
  }

  @Action({ commit: "setUser" })
  async signup(userCreateReq: UserCreate) {
    try {
      await createUser(userCreateReq)
    
      const wallet = hedgehog.getWallet();
      const privKeyBuf = wallet.getPrivateKey();
      const identity: HedgehogIdentity = new HedgehogIdentity(privKeyBuf);
      this.context.commit("setHedgehogIdent", identity);
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
  async getLoggedInUser() {
    try {
      let identity: HedgehogIdentity | null = this.hedgehogIdent;
      if (!identity) {
        const wallet = hedgehog.getWallet();
        const privKeyBuf = wallet.getPrivateKey();
        identity = new HedgehogIdentity(privKeyBuf);
      }
      this.context.commit("setHedgehogIdent", identity);
      const respData: any = await getUserFn(identity);
      console.log(respData);

      const user: User = {
        id: respData.id,
        displayName: respData.displayName,
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
      if (!this.hedgehogIdent) {
        throw new Error("Hedgehog identity not set");
      }
      if (!this.user) {
        throw new Error("no logged in user")
      }
      const respData: User = await updateUserFn(this.user.id, userUpdateReq, this.hedgehogIdent);
      
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

  @Mutation
  setTextileIdent(identity: TextileIdentity) {
    this.textileIdent = identity;
  }

  @Action
  async setupUserBuckets() {
    try {
      const wallet = hedgehog.getWallet();
      console.log(wallet.getPublicKeyString());

      const privKeyBuf = wallet.getPrivateKey();
      const key = await keys.supportedKeys.secp256k1.unmarshalSecp256k1PrivateKey(
        privKeyBuf
      );

      const identity: TextileIdentity = new TextileIdentity(key);
      this.context.commit("setTextileIdent", identity);

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
      if (!this.hedgehogIdent) {
        throw new Error("Hedgehog identity not set");
      }
      const resp = await uploadTrackToBucket(
        recordedTrack,
        this.bucketKey,
        this.buckets
      );

      const trackMetadata: StoreTrackMetadata = {
        cid: resp.cid,
        title: resp.name,
        parentTrackID: recordedTrack.parentTrackID,
        instruments: Array.from(recordedTrack.instrumentTags)
      }

      const respData: any = await storeTrackFn(trackMetadata, this.hedgehogIdent);

      const track: TrackMetadata = {
        id: respData.id,
        cid: respData.cid,
        title: respData.title,
        composerID: respData.composerID,
        parentTrackID: respData.parentTrackID,
        nForks: respData.nForks,
        instruments: respData.instruments
      };
      this.context.commit("updateTracks", track);
    } catch (err) {
      console.error(err);
    }
  }

  @Mutation
  setUserFeed(userFeed: UserFeed) {
    this.userFeed = userFeed;
  }

  @Action
  async loadUserFeed() {
    try {
      if (!this.hedgehogIdent) {
        throw new Error("Hedgehog identity not set");
      }
      if (!this.user) {
        throw new Error("no logged in user")
      }

      const respData: Array<object> = await getUserFeedFn(this.hedgehogIdent);
      const userFeed: UserFeed = {
        tracks: []
      }

      for (let i = 0; i < respData.length; i++) {
        const feedElem: any = respData[i];
        userFeed.tracks = [
          ...userFeed.tracks,
          {
            name: feedElem.Name,
            tracks: feedElem.Tracks
          }
        ];
      }

      this.context.commit("setUserFeed", userFeed);
    } catch (err) {
      console.error(err)
    }
  }
}

export default getModule(UsersModule);
