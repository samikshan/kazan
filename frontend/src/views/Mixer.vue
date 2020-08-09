<template>
  <v-app>
    <v-skeleton-loader
      :loading="!readyToMix"
      type="card"
    >
      <v-card
        class="pa-2 ma-2 d-flex align-center justify-space-around mb-6"
        dark
        color="grey darken-1"
      >
        <v-btn icon v-if="readyToMix" @click="handleMixButtonClick">
          <v-icon x-large dark>{{ isMixing ? "mdi-stop" : "mdi-microphone" }}</v-icon>
        </v-btn>
      </v-card>
      <v-card
        v-if="recordings.length > 0"
        dark
        color="grey darken-1"
      >
        <v-card
          class="pa-2 ma-2 d-flex align-center justify-space-around mb-6"
          color="grey darken-1"
        >
          <v-card-title
            class=".text-lg-h6"
          >
            Mixes
          </v-card-title>
        </v-card>
        <RecordingList
          :recordings="recordings"
          @add-instrument-tag="handleAddTag($event)"
          @remove-instrument-tag="handleRemoveTag($event)"
          @publish-recording="handleUploadRecording($event)"
        ></RecordingList>
      </v-card>
    </v-skeleton-loader>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import axios from "axios";
import RecordingList from "@/components/RecordingList.vue";
import { RecordedTrack } from "@/store/models"
import users from "@/store/modules/users";

@Component({
  components: {
    RecordingList
  }
})
export default class Mixer extends Vue {
  audioCtx: AudioContext = new AudioContext({
    latencyHint: 0
  });
  mediaStreamSrc!: MediaStreamAudioSourceNode;
  mediaElemSrc!: MediaElementAudioSourceNode;
  mixedAudio!: MediaStreamAudioDestinationNode;
  parentAudioElem!: HTMLMediaElement;

readyToMix = false;
  isMixing = false;
  chunks: Array<any> = []
  mediaRecorder!: MediaRecorder;
  recordings: Array<RecordedTrack> = [];

  @Prop() parentAudioCID!: string;

  get parentAudioSrc() {
    return `https://${this.$route.params.parentCID}.ipfs.hub.textile.io`
  }

  async created() {
    try {
      this.parentAudioElem = new Audio(this.parentAudioSrc);
      this.parentAudioElem.crossOrigin = "anonymous";
      this.parentAudioElem.addEventListener("canplaythrough", event => {
        this.readyToMix = true
      });

      this.mixedAudio = this.audioCtx.createMediaStreamDestination();
      this.mediaElemSrc = this.audioCtx.createMediaElementSource(this.parentAudioElem);
      this.mediaElemSrc.connect(this.audioCtx.destination);
      this.mediaElemSrc.connect(this.mixedAudio);
    } catch(err) {
      console.error(err);
    }
  }

  async handleMixButtonClick() {
    if (this.isMixing) {
      await this.stopMixing()
    } else {
      await this.startMixing()
    }
  }

  async startMixing() {
    console.log("started mixing");
    this.isMixing = true;
    let stream = null;
    try {
      stream = await navigator.mediaDevices.getUserMedia({
        audio: true,
        video: false
      });

      this.mediaStreamSrc = this.audioCtx.createMediaStreamSource(stream);
      this.mediaStreamSrc.connect(this.mixedAudio);
      this.mediaStreamSrc.connect(this.audioCtx.destination);
      this.parentAudioElem.play();

      /* use the stream */
      this.mediaRecorder = new MediaRecorder(this.mixedAudio.stream);
      this.mediaRecorder.start(1);
      this.mediaRecorder.ondataavailable = event =>
        this.chunks.push(event.data);

      this.mediaRecorder.onstop = () => this.handleMediaRecorderStop();
    } catch (err) {
      console.log("The following error occurred: " + err);
    }
  }

  stopMixing() {
    console.log("stopped mixing");
    this.isMixing = false;
    this.mediaRecorder.stop();
    this.mediaStreamSrc.disconnect(this.audioCtx.destination);
  }

  handleMediaRecorderStop() {
    console.log("data available after MediaRecorder.stop() called.");
    let clipName = prompt("Enter a name for your sound clip");
    if (!clipName) {
      console.error("failed to get clipname");
      clipName = "";
    }

    let blobOpts = { type: "audio/wav; codecs=0" };
    if (MediaRecorder.isTypeSupported("audio/wav;codecs=MS_PCM")) {
      blobOpts = { type: "audio/wav; codecs=MS_PCM" };
    }

    const blob = new Blob(this.chunks, blobOpts);
    this.chunks = [];
    const audioURL = URL.createObjectURL(blob);
    console.log("recorder stopped");
    console.log(audioURL);

    this.recordings = [
      ...this.recordings,
      {
        id: this.recordings.length,
        data: blob,
        localURL: audioURL,
        name: clipName,
        instrumentTags: new Set<string>(),
        isPublished: false
      }
    ];
  }

  handleAddTag(event: {tagText: string, recordingID: number}) {
    console.log(event);
    const recID = event.recordingID;
    this.recordings[recID].instrumentTags.add(event.tagText);
  }

  handleRemoveTag(event: {tagText: string, recordingID: number}) {
    console.log(event);
    const recID = event.recordingID;
    this.recordings[recID].instrumentTags.delete(event.tagText);
  }

  async handleUploadRecording(recID: number) {
    console.log("upload request for recording id: ", recID);
    const recordedTrack = this.recordings[recID];
    try {
      await users.addNewTrack(recordedTrack);
      recordedTrack.isPublished = true;
    } catch (e) {
      console.error(e);
    }
  }

  handleDeleteRecording(event: Event) {
    const target = event.target as HTMLButtonElement;
    console.log("delete request for recording id: ", target.id);
  }
}

</script>