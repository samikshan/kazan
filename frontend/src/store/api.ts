import { Libp2pCryptoIdentity } from "@textile/threads-core";

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
