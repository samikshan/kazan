<template>
  <v-app>
    <v-card> 
      <v-btn icon v-if="readyToMix" @click="handleMixButtonClick">
        <v-icon color="red">{{ isMixing ? "mdi-stop" : "mdi-microphone" }}</v-icon>
      </v-btn>

      <v-card class="mx-auto">
        <v-card-title>
          <h2 class="display-1">Mixes</h2>
        </v-card-title>
        <v-card>
          <ul class="audio-rec-list">
            <li v-for="r in recordings" v-bind:key="r.id">
              <p>{{ r.name }}</p>
              <audio controls v-bind:src="r.localURL"></audio>
            </li>
          </ul>
        </v-card>
      </v-card>
    </v-card>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import axios from "axios";
import { RecordedTrack } from "@/store/models"

@Component
export default class Mixer extends Vue {
  audioCtx: AudioContext = new AudioContext();
  mediaStreamSrc!: MediaStreamAudioSourceNode;
  mediaElemSrc!: MediaElementAudioSourceNode;
  mixedAudio!: MediaStreamAudioDestinationNode;
  merger!: ChannelMergerNode;
  splitter!: ChannelSplitterNode;
  parentAudioElem!: HTMLMediaElement;
  channel1 = [0,1];
  channel2 = [1, 0];

  readyToMix = false;
  isMixing = false;
  chunks: Array<any> = []
  mediaRecorder!: MediaRecorder;
  recordings: Array<RecordedTrack> = [];

  @Prop() parentAudioCID!: string;

  get parentAudioSrc() {
    return `https://${this.parentAudioCID}.ipfs.hub.textile.io`
  }

  async created() {
    // let stream = null;
    try {
      this.parentAudioElem = new Audio(this.parentAudioSrc);
      this.parentAudioElem.crossOrigin = "anonymous";
      this.parentAudioElem.addEventListener("canplaythrough", event => {
        this.readyToMix = true
      });

      this.merger = this.audioCtx.createChannelMerger(2);
      this.splitter = this.audioCtx.createChannelSplitter(2);
      this.mixedAudio = this.audioCtx.createMediaStreamDestination();
      this.mediaElemSrc = this.audioCtx.createMediaElementSource(this.parentAudioElem);
      this.mediaElemSrc.connect(this.splitter);
      console.log("media element source connected to splitter");
      this.splitter.connect(this.merger, this.channel1[0], this.channel1[1]);

      console.log("media element source connected to merger");

      this.merger.connect(this.mixedAudio);
      this.merger.connect(this.audioCtx.destination);
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
      this.mediaStreamSrc.connect(this.splitter);
      console.log("stream source connected to splitter");
      this.splitter.connect(this.merger, this.channel2[0], this.channel2[1]);
      console.log("stream source connected to merger");
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
    this.mediaStreamSrc.disconnect(this.splitter);
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
}

</script>