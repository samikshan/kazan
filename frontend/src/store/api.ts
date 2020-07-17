import { Libp2pCryptoIdentity } from "@textile/threads-core";

export async function getIdentity(): Promise<Libp2pCryptoIdentity> {
	try {
		const storedIdent = localStorage.getItem("identity")
		if (storedIdent === null) {
			throw new Error('No identity')
		}
		const restored = Libp2pCryptoIdentity.fromString(storedIdent)
		return restored
	}
	catch (e) {
		/**
		 * If any error, create a new identity.
		 */
		try {
			const identity = await Libp2pCryptoIdentity.fromRandom()
			const identityString = identity.toString()
			localStorage.setItem("identity", identityString)
			return identity
		} catch (err) {
			return err.message
		}
	}
}