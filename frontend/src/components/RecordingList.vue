<template>
  <v-container>
    <v-row dense>
      <v-col
        v-for="(recording, index) in recordings"
        :key="index"
        cols="12"
      >
        <v-card>
          <Recording :title="recording.name" :localURL="recording.localURL" />
          <v-card-actions v-if="!recording.isPublished">
            <v-dialog v-model="dialogOpen" persistent max-width="600px">
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  v-bind="attrs"
                  v-on="on"
                >
                  <v-icon>mdi-upload</v-icon>
                  Add Tags And Publish
                </v-btn>
              </template>
              <v-card>
                <v-card-title>
                  <span class="headline">Add Track Information</span>
                </v-card-title>
                <TrackInfoInput @instrument-tag-added="instrumentTagAdded($event, index)" />
                <v-card-actions>
                <v-spacer></v-spacer>
                  <v-btn
                    color="blue darken-1"
                    text
                    @click="publish(index)"
                  >Publish</v-btn>
                </v-card-actions>
              </v-card>
            </v-dialog>
            <v-btn>
              <v-icon>mdi-delete</v-icon>
              Reject 
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import Recording from "@/components/Recording.vue";
import TrackInfoInput from "@/components/TrackInfoInput.vue";
import { RecordedTrack } from "@/store/models";
import users from "@/store/modules/users";

@Component({
  components: {
    Recording,
    TrackInfoInput
  }
})
export default class RecordingList extends Vue {
  dialogOpen = false;
  @Prop() recordings!: Array<RecordedTrack>

  instrumentTagAdded(tag: string, index: number) {
    const tagInfo = {
      tagText: tag,
      recordingID: index
    }
    this.$emit("add-instrumentt-tag", tagInfo);
  }

  instrumentTagRemoved(tag: string, index: number) {
    const tagInfo = {
      tagText: tag,
      recordingID: index
    }
    this.$emit("remove-instrument-tag", tagInfo);
  }

  publish(recordingID: number) {
    this.dialogOpen = false;
    this.$emit("publish-recording", recordingID);
  }
}
</script>