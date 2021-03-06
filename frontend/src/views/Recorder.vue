<template>
  <v-app>
    <v-card class="mx-auto mt-5" max-width="600px" min-width="360px">
      <v-card-actions elevation=24>
        <div id="buttons">
          <v-btn @click="handleStartRecording">
            <v-icon large color="teal lighten-1">{{ isRecording ? 'mdi-pause' : 'mdi-microphone' }}</v-icon>
            {{ isRecording ? "Pause Recording" : "Start Recording" }}
          </v-btn>
          <v-btn @click="handleStopRecording">
            <v-icon large color="red">mdi-stop</v-icon>
            Stop Recording
          </v-btn>
        </div>
      </v-card-actions>

      <v-card class="mx-auto">
        <v-card-title>
          <h6>Recorded tracks</h6>
        </v-card-title>
        <v-card>
          <RecordingList
            :recordings="recordings"
            @add-instrument-tag="handleAddTag($event)"
            @remove-instrument-tag="handleRemoveTag($event)"
            @publish-recording="handleUploadRecording($event)"
          ></RecordingList>
        </v-card>
      </v-card>
    </v-card>
  </v-app>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import RecordingList from "@/components/RecordingList.vue";
import { RecordedTrack } from "@/store/models";
import users from "@/store/modules/users";

@Component({
  components: {
    RecordingList
  }
})
export default class Recorder extends Vue {
  audioCtx: AudioContext = new AudioContext({
    latencyHint: 0
  });
  mediaStreamSrc!: MediaStreamAudioSourceNode;
  recordedAudio!: MediaStreamAudioDestinationNode;
  isRecording = false;
  recordings: Array<RecordedTrack> = [];
  mediaRecorder!: MediaRecorder;
  chunks: Array<any> = [];

  async mounted() {
    try {
      console.log(this.audioCtx.baseLatency);
      this.recordedAudio = this.audioCtx.createMediaStreamDestination();
    } catch (err) {
      /* handle the error */
      console.log("The following error occurred: " + err);
    }
  }

  async handleStartRecording(event: Event) {
    console.log("started recording");
    this.isRecording = true;
    let stream = null;
    try {
      stream = await navigator.mediaDevices.getUserMedia({
        audio: true,
        video: false
      });

      this.mediaStreamSrc = this.audioCtx.createMediaStreamSource(stream);
      this.mediaStreamSrc.connect(this.audioCtx.destination);

      /* use the stream */
      this.mediaRecorder = new MediaRecorder(this.recordedAudio.stream);
      this.mediaRecorder.start(1);
      this.mediaRecorder.ondataavailable = event =>
        this.chunks.push(event.data);

      this.mediaRecorder.onstop = () => this.handleMediaRecorderStop();
    } catch (err) {
      console.error(err);
    }
  }

  stopRecording() {
    console.log("stopped recording");
    this.isRecording = false;
    this.mediaRecorder.stop();
  }

  handleStopRecording(event: Event) {
    console.log(this.mediaRecorder.state);
    this.isRecording = false;
    this.mediaRecorder.stop();
    this.mediaStreamSrc.disconnect(this.audioCtx.destination);
    console.log(this.mediaRecorder.state);
    console.log("recording stopped");
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
    // const target = event.target as HTMLButtonElement;
    console.log("upload request for recording id: ", recID);
    // const recID: number = +target.id;
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
