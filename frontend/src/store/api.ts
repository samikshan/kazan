import { Buckets, Identity, KeyInfo, Client} from '@textile/hub';
import { Hedgehog } from "@audius/hedgehog";
import axios from "axios";
import { UserCreate } from "@/store/models";

const requestToServer = async (axiosRequestObj: any) => {
  axiosRequestObj.baseURL = 'http://localhost:1323/'

  try {
    const resp = await axios(axiosRequestObj)
    if (resp.status === 200) {
      return resp.data
    } else {
      throw new Error('Server returned error: ' + resp.status.toString() + ' ' + resp.data['error'])
    }
  } catch (e) {
    throw new Error('Server returned error: ' + e.response.status.toString() + ' ' + e.response.data['error'])
  }
}

const setAuthFn = async (obj: any) => {
  await requestToServer({
    url: '/authentication',
    method: 'post',
    data: obj
  })
}

const setUserFn = async (obj: any) => {
  await requestToServer({
    url: '/user',
    method: 'post',
    data: obj
  })
}

const getFn = async (obj: any) => {
  return requestToServer({
    url: '/authentication',
    method: 'get',
    params: obj
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
}

export async function createThreadsClient(identity: Identity) {
  const client = await Client.withKeyInfo(keyinfo);
  const token = await client.getToken(identity);
  console.log(token)
  return client;
}

export async function createBucketsClient(identity: Identity) {
  const buckets = await Buckets.withKeyInfo(keyinfo);
  console.log(buckets)
  const token = await buckets.getToken(identity);
  return buckets;
}

export async function createBucket(buckets: Buckets, name: string) {
  const root = await buckets.open(name)
  if (!root) {
    throw new Error('Failed to open bucket')
  }

  return root.key
}