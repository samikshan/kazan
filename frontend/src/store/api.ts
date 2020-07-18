import { Libp2pCryptoIdentity } from "@textile/threads-core";
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

export function getStoredIdentity(): Promise<Libp2pCryptoIdentity> {
  try {
		const storedIdent = localStorage.getItem("identity")
		if (storedIdent === null) {
			throw new Error('No identity')
		}
		const restored = Libp2pCryptoIdentity.fromString(storedIdent)
		return restored
	} catch (e) {
    return e.message;
  }
}

export async function createIdentity(): Promise<Libp2pCryptoIdentity> {
	try {
    const identity = await Libp2pCryptoIdentity.fromRandom()
    const identityString = identity.toString()
    localStorage.setItem("identity", identityString)
    return identity
  } catch (err) {
    return err.message
  }
}

export async function createUser(user: UserCreate) {
  const wallet = await hedgehog.signUp(user.username, user.password);
  console.log(wallet);
}