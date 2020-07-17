import { Identity } from "@textile/threads-core";

export interface Profile {
  id: number
  username: string
  fullName?: string
  image?: string
}

export interface User {
  identity: Identity 
}

export interface Track {
  cid: string
}

export interface RecordedTrack {
  id: number
  localURL: string
  data: Blob
}