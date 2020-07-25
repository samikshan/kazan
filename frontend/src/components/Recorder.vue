<template>
  <div class="record-upload-track">
    <section class="main-controls">
      <canvas class="visualizer" height="60px"></canvas>

      <div id="buttons">
        <button
          class="record-start"
          id="startRecording"
          v-on:click="handleStartRecording"
        >
          Start recording
        </button>
        <button
          class="record-stop"
          id="stopRecording"
          v-on:click="handleStopRecording"
        >
          Stop recording
        </button>
      </div>
    </section>

    <hr />

    <h2>Recordings</h2>

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
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { RecordedTrack } from "@/store/models";
import users from "@/store/modules/users";

@Component
export default class Recorder extends Vue {
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
    const target = event.target as HTMLButtonElement;
    target.disabled = true;
    const stopRecording = document.querySelector(
      ".record-stop"
    ) as HTMLButtonElement;
    stopRecording.disabled = false;
    console.log(this.mediaRecorder);
    this.mediaRecorder.start();
    console.log(this.mediaRecorder.state);
    console.log("recording started");
  }

  handleStopRecording(event: Event) {
    console.log(this.mediaRecorder.state);
    const startRecording = document.querySelector(
      ".record-start"
    ) as HTMLButtonElement;
    startRecording.disabled = false;
    const target = event.target as HTMLButtonElement;
    target.disabled = true;
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
    // recordedBlob.name = 'track_' + recID + '.wav';
    // recordedBlob.lastModifiedDate = new Date();
    // console.log(recordedBlob);
    // const cid = await hotUpload(recordedBlob);
    // const jobID = await coldUpload(cid);
    // console.log(jobID);

    // try {
    //     const newTrackData = { title: "", cID: cid };
    //     await postNewTrack(newTrackData);
    // } catch (err) {
    //     console.error("upload failed. Returned error: ", err);
    //     await showAlert("Uploading new track failed!");
    //     return;
    // }
    // var recordingID = target.id;
    // this.recordings = this.recordings.filter(recording => recording.id != recordingID);
  }

  handleDeleteRecording(event: Event) {
    const target = event.target as HTMLButtonElement;
    console.log("delete request for recording id: ", target.id);
    // var recordingID = target.id;
    // this.recordings = this.recordings.filter(recording => recording.id != recordingID);
  }
}
</script>
