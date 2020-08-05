export interface User {
  id: number;
  username: string;
  walletAddr: string;
  instruments: Array<Instrument>;
}

export interface Instrument {
  name: string
}

export interface UserCreate {
  username: string;
  password: string;
}

export interface UserCreateResponse {
  user: User;
}

export interface UserUpdate {
  displayName?: string;
  instruments?: Array<string>;
}

export interface UserUpdateResponse {
  user: User;
}

export interface BucketUploadResponse {
  name: string;
  cid: string;
}

export interface TrackUploadedData {
  name: string;
  cid: string;
  path: string;
}

export interface TrackData {
  name: string;
  date: number;
  metadata: TrackUploadedData;
}

export interface Track {
  name: string;
  cid: string;
  src: string;
}

export interface StoreTrackMetadata {
  cid: string;
  title: string;
  parentTrackID?: number;
  components: Array<string>;
}

export interface TrackMetadata {
  cid: string;
	title: string;
	composerID: number;
	composer: User;
	parentTrackID: number;
	parentTrack: TrackMetadata;
	forks: Array<TrackMetadata>;
	components: Array<Component>;
}

export interface Component {
  name: string;
  tracks: Array<TrackMetadata>;
}

export interface StoreTrackMetadataResp {
  cid: string;
	title: string;
	composerID: number;
	composer: User;
	parentTrackID: number;
	parentTrack: TrackMetadata;
	forks: Array<TrackMetadata>;
	components: Array<Component>;
}

export interface UserTrackIndex {
  owner: string;
  date: number;
  paths: string[];
}

export interface RecordedTrack {
  id: number;
  localURL: string;
  data: Blob;
  name: string;
  instrumentTags: Set<string>;
  isPublished: boolean;
  parentTrackID?: number;
}
