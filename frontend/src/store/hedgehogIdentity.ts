import { ecsign, ECDSASignature, privateToPublic, ecrecover, bufferToHex, toRpcSig, toBuffer, fromRpcSig } from "ethereumjs-util";
import { Public, Identity } from "@textile/threads-core";

export class HedgehogPublicKey implements Public {
  constructor(public key: Buffer) {}

  /**
   * Verify the given signed data.
   * @param msgHash The message hash to verify.
   * @param sig The signature to verify.
   */
  verify(msgHash: Uint8Array, sig: Uint8Array): Promise<boolean> {
    return new Promise((resolve, reject) => {
      try {
        const sigBuf: Buffer = Buffer.from(sig);
        const sigHex: string = bufferToHex(sigBuf);
        const ecdsaSig: ECDSASignature = fromRpcSig(sigHex);
        const msgHashBuf: Buffer = Buffer.from(msgHash);
        const recoveredKey: Buffer = ecrecover(msgHashBuf, ecdsaSig.v, ecdsaSig.r, ecdsaSig.s);
        resolve(recoveredKey == this.key);
      } catch (err) {
        console.error(err);
        reject(err);
      }
    });
  }

  /**
   * Returns 0x prefixed hex string
   */
  toString(): string {
    return bufferToHex(this.key);
  }

  /**
   * The raw bytes of the Public key.
   */
  get bytes(): Uint8Array {
    const keyBytes: Uint8Array = new Uint8Array(this.key);
    return keyBytes;
  }
}

export class HedgehogIdentity implements Identity {
  key: Buffer;
  constructor(key: Buffer) {
    this.key = key;
  }

  /**
   * Signs the given hash message with the Private key,
   * @param msgHash Hash of the message to be signed
   */
  sign(msgHash: Uint8Array): Promise<Uint8Array> {
    return new Promise((resolve, reject) => {
      try {
        const msgHashBuf: Buffer = Buffer.from(msgHash);
        const ecdsaSig: ECDSASignature = ecsign(msgHashBuf, this.key);
        const sigHex: string = toRpcSig(ecdsaSig.v, ecdsaSig.r, ecdsaSig.s);
        const sigBuf: Buffer = toBuffer(sigHex);
        const sig: Uint8Array = new Uint8Array(sigBuf);
        resolve(sig);
      } catch (err) {
        console.error(err);
        reject(err);
      }
    });
  }

  /**
   * Returns the Public key.
   */
  get public(): HedgehogPublicKey {
    return new HedgehogPublicKey(privateToPublic(this.key));
  }

  /**
   * Returns base32 encoded private key representation.
   */
  // toString(): string {
  //   return privateKeyToString(this.key);
  // }
}
