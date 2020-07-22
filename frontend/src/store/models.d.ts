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

export interface UploadTrackResponse {
  cid: string
  metapath: string
  name: string
}

export interface TrackUploadedData {
  name: string
  cid: string
  path: string
}

export interface TrackData {
  name: string
  date: number
  metadata: TrackUploadedData
}

export interface Track {
  name: string
  cid: string
  src: string
}

export interface UserTrackIndex {
  owner: string
  date: number
  paths: string[]
}

export interface RecordedTrack {
  id: number
  localURL: string
  data: Blob
  name: string
}
