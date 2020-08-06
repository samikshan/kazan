import { Buckets, Identity, Client } from "@textile/hub";
import { Hedgehog } from "@audius/hedgehog";
import axios from "axios";
import {
  UserCreate,
  RecordedTrack,
  BucketUploadResponse,
  UserTrackIndex,
  Track,
  TrackData,
  StoreTrackMetadata,
  UserUpdate
} from "@/store/models";
import { bufferToHex, hashPersonalMessage, ecrecover, ECDSASignature, fromRpcSig } from 'ethereumjs-util';
import { HedgehogIdentity } from './hedgehogIdentity';

const requestToServer = async (axiosRequestObj: any) => {
  axiosRequestObj.baseURL = "http://localhost:1323/";

  try {
    const resp = await axios(axiosRequestObj);
    if (resp.status === 200) {
      return resp.data;
    } else {
      throw new Error(
        "Server returned error: " +
          resp.status.toString() +
          " " +
          resp.data["error"]
      );
    }
  } catch (e) {
    throw new Error(
      "Server returned error: " +
        e.response.status.toString() +
        " " +
        e.response.data["error"]
    );
  }
};

const setAuthFn = async (obj: any) => {
  await requestToServer({
    url: "/authentication",
    method: "post",
    data: obj
  });
};

const setUserFn = async (obj: any) => {
  await requestToServer({
    url: "/user",
    method: "post",
    data: obj
  });
};

const getFn = async (obj: any) => {
  return requestToServer({
    url: "/authentication",
    method: "get",
    params: obj
  });
};

export const storeTrackFn = async (obj: StoreTrackMetadata, identity: Identity) => {
  const msg = "Authenticate with Kazan service";
  const msgBuf = Buffer.from(msg);
  const msgHashBuf: Buffer = hashPersonalMessage(msgBuf);
  const msgHash: Uint8Array = new Uint8Array(msgHashBuf);
  const sig = await identity.sign(msgHash);
  const sigBuf: Buffer = Buffer.from(sig);
  const sigHex: string = bufferToHex(sigBuf);
  const msgHashHex: string = bufferToHex(msgHashBuf);
  return requestToServer({
    url: "/tracks",
    method: "post",
    data: obj,
    headers: {
      'encoded-data-message': msgHashHex,
      'encoded-data-signature': sigHex,
      'Content-Type': 'application/json'
    }
  });
}

export const getUserFn = async (identity: Identity) => {
  const msg = "Authenticate with Kazan service";
  const msgBuf = Buffer.from(msg);
  const msgHashBuf: Buffer = hashPersonalMessage(msgBuf);
  const msgHash: Uint8Array = new Uint8Array(msgHashBuf);
  const sig = await identity.sign(msgHash);
  const sigBuf: Buffer = Buffer.from(sig);
  const sigHex: string = bufferToHex(sigBuf);
  const msgHashHex: string = bufferToHex(msgHashBuf);
  
  return requestToServer({
    url: "/user",
    method: "get",
    headers: {
      'encoded-data-message': msgHashHex,
      'encoded-data-signature': sigHex,
      'Content-Type': 'application/json'
    }
  })
}

export const getUserFeedFn = async (identity: HedgehogIdentity) => {
  const msg = "Authenticate with Kazan service";
  const msgBuf = Buffer.from(msg);
  const msgHashBuf: Buffer = hashPersonalMessage(msgBuf);
  const msgHash: Uint8Array = new Uint8Array(msgHashBuf);
  const sig = await identity.sign(msgHash);
  const sigBuf: Buffer = Buffer.from(sig);
  const sigHex: string = bufferToHex(sigBuf);
  const msgHashHex: string = bufferToHex(msgHashBuf);
  
  return requestToServer({
    url: "/feed",
    method: "get",
    headers: {
      'encoded-data-message': msgHashHex,
      'encoded-data-signature': sigHex,
      'Content-Type': 'application/json'
    }
  });
}

export const updateUserFn = async (userID: number, obj: UserUpdate, identity: Identity) => {
  // const challenge = Buffer.from('Sign this string');
  const msg = "Authenticate with Kazan service";
  const msgBuf = Buffer.from(msg);
  const msgHashBuf: Buffer = hashPersonalMessage(msgBuf);
  const msgHash: Uint8Array = new Uint8Array(msgHashBuf);
  const sig = await identity.sign(msgHash);
  const sigBuf: Buffer = Buffer.from(sig);
  const sigHex: string = bufferToHex(sigBuf);
  const sigEcdsa : ECDSASignature = fromRpcSig(sigHex);
  const msgHashHex: string = bufferToHex(msgHashBuf);
  console.log(sigHex);
  console.log(msgHashHex);

  const pubkey: string = bufferToHex(ecrecover(msgHashBuf, sigEcdsa.v, sigEcdsa.r, sigEcdsa.s));
  console.log(pubkey);
  console.log(identity.public.toString())

  return requestToServer({
    url: "/user/" + userID,
    method: "put",
    headers: {
      'encoded-data-message': msgHashHex,
      'encoded-data-signature': sigHex,
      'Content-Type': 'application/json'
    },
    data: obj
  })
}

export const hedgehog = new Hedgehog(getFn, setAuthFn, setUserFn);

export async function createUser(user: UserCreate) {
  const wallet = await hedgehog.signUp(user.username, user.password);
  console.log(wallet);
}

// textile stuff

const keyinfo = {
  key: "bfn2ssz72cgn6cgtuayqwnenvgi",
  secret: "",
  type: 1
};

export async function createThreadsClient(identity: Identity) {
  const client = await Client.withKeyInfo(keyinfo);
  const token = await client.getToken(identity);
  console.log(token);
  return client;
}

export async function createBucketsClient(identity: Identity) {
  const buckets = await Buckets.withKeyInfo(keyinfo);
  console.log(buckets);
  const token = await buckets.getToken(identity);
  console.log(token);
  return buckets;
}

export async function createBucket(buckets: Buckets, name: string) {
  const root = await buckets.open(name);
  if (!root) {
    throw new Error("Failed to open bucket");
  }

  return root.key;
}

export async function uploadTrackToBucket(
  recordedTrack: RecordedTrack,
  bucketKey: string,
  buckets: Buckets
): Promise<BucketUploadResponse> {
  return new Promise((resolve, reject) => {
    try {
      const now = new Date().getTime();
      const filename = `${now}_${recordedTrack.name}`;

      const data = recordedTrack.data;
      data.arrayBuffer().then(audioBuffer => {
        const path = `tracks/${filename}`;
        console.log("pushing to bucket path: ", path);
        buckets.pushPath(bucketKey, path, audioBuffer).then(raw => {
          const response: BucketUploadResponse = {
            name: filename,
            cid: raw.path.cid.toString()
          };

          resolve(response);
        });
      });
    } catch (e) {
      reject(e);
    }
  });
}

export async function initIndex(identity: Identity): Promise<UserTrackIndex> {
  const index: UserTrackIndex = {
    owner: identity.public.toString(),
    date: new Date().getTime(),
    paths: []
  };
  return index;
}

export async function storeIndex(
  index: UserTrackIndex,
  buckets: Buckets,
  bucketKey: string
) {
  const buf = Buffer.from(JSON.stringify(index, null, 2));
  const path = "index.json";
  await buckets.pushPath(bucketKey, path, buf);
}

export async function getTrackIndex(
  buckets: Buckets,
  bucketKey: string,
  identity: Identity
) {
  try {
    const metadata = buckets.pullPath(bucketKey, "index.json");
    const { value } = await metadata.next();
    let str = "";
    for (let i = 0; i < value.length; i++) {
      str += String.fromCharCode(parseInt(value[i]));
    }
    const index: UserTrackIndex = JSON.parse(str);
    return index;
  } catch (error) {
    const index = await initIndex(identity);
    return index;
  }
}

export async function getTracks(
  buckets: Buckets,
  bucketKey: string,
  index: UserTrackIndex
): Promise<Track[]> {
  return new Promise((resolve, reject) => {
    try {
      const tracks: Track[] = [];
      for (const path of index.paths) {
        const metadata = buckets.pullPath(bucketKey, path);
        metadata.next().then(iterResult => {
          const value = iterResult.value;
          let str = "";
          for (let i = 0; i < value.length; i++) {
            str += String.fromCharCode(parseInt(value[i]));
          }
          const trackData: TrackData = JSON.parse(str);
          const track: Track = {
            src: `https://${trackData.metadata.cid}.ipfs.hub.textile.io`,
            name: trackData.name,
            cid: trackData.metadata.cid
          };
          tracks.push(track);
        });
      }
      resolve(tracks);
    } catch (err) {
      reject(err);
    }
  });
}
