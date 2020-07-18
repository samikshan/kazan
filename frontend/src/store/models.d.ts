import { Identity } from "@textile/threads-core";

export interface User {
  username: string
  walletAddr: string
  email?: string
}

export interface UserCreate {
  username: string;
  password: string;
}

export interface UserCreateResponse {
  user: User
}

export interface Track {
  cid: string
}

export interface RecordedTrack {
  id: number
  localURL: string
  data: Blob
}