<template>
  <v-app>
    <v-card class="mx-auto mt-5" max-width="600px" min-width="360px">
      <v-card-actions elevation=24>
        <div id="buttons">
          <v-btn @click="handleStartRecording">
            <v-icon color="teal lighten-1">{{ isRecording ? 'mdi-pause' : 'mdi-microphone' }}</v-icon>
          </v-btn>
          <v-btn @click="handleStopRecording">
            <v-icon color="red">mdi-stop</v-icon>
          </v-btn>
        </div>
      </v-card-actions>

      <v-card class="mx-auto">
        <v-card-title>
          <h2 class="display-1">Recorded tracks</h2>
        </v-card-title>
        <v-card>
          <ul class="audio-rec-list">
            <li v-for="r in recordings" v-bind:key="r.id">
              <p>{{ r.name }}</p>
              <audio controls v-bind:src="r.localURL"></audio>
              <div id="buttons">
                <button v-bind:id="r.id" v-on:click="handleUploadRecording">
                  Publish
                </button>
                <button v-bind:id="r.id" v-on:click="handleDeleteRecording">
                  Delete
                </button>
              </div>
            </li>
          </ul>
        </v-card>
      </v-card>
    </v-card>
  </v-app>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { RecordedTrack } from "@/store/models";
import users from "@/store/modules/users";

@Component
export default class Recorder extends Vue {
  isRecording = false;
  recordings: Array<RecordedTrack> = [];
  mediaRecorder!: MediaRecorder;
  chunks: Array<any> = [];

  async mounted() {
    let stream = null;
    try {
      stream = await navigator.mediaDevices.getUserMedia({
        audio: true,
        video: false
      });
      /* use the stream */
      this.mediaRecorder = new MediaRecorder(stream, {
        mimeType: "audio/webm"
      });
      this.mediaRecorder.ondataavailable = event =>
        this.chunks.push(event.data);
      this.mediaRecorder.onstop = event => this.handleMediaRecorderStop(event);
    } catch (err) {
      /* handle the error */
      console.log("The following error occurred: " + err);
    }
  }

  handleStartRecording(event: Event) {
    console.log(this.mediaRecorder);
    this.mediaRecorder.start();
    console.log(this.mediaRecorder.state);
    console.log("recording started");
  }

  handleStopRecording(event: Event) {
    console.log(this.mediaRecorder.state);
    this.mediaRecorder.stop();
    console.log(this.mediaRecorder.state);
    console.log("recording stopped");
  }

  handleMediaRecorderStop(event: Event) {
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
    const audioURL = URL.createObjectURL(blob);
    console.log("recorder stopped");
    console.log(audioURL);

    this.recordings = [
      ...this.recordings,
      {
        id: this.recordings.length,
        data: blob,
        localURL: audioURL,
        name: clipName
      }
    ];
  }

  async handleUploadRecording(event: Event) {
    const target = event.target as HTMLButtonElement;
    console.log("upload request for recording id: ", target.id);
    const recID: number = +target.id;
    const recordedTrack = this.recordings[recID];
    try {
      await users.addNewTrack(recordedTrack);
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
