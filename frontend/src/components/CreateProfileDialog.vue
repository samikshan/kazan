<template>
  <v-card class="mx-auto mt-5" max-width="600px" min-width="360px">
    <v-card v-if="!createProfile">
      <v-card-title class="headline">You're signed up!</v-card-title>
      <v-card-text>Your wallet address is: {{ address }}</v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="green darken-1" text @click="createProfile = true">Continue to update profile details</v-btn>
      </v-card-actions>
    </v-card>
    <v-card v-else>
      <v-card-title class="headline">Update Profile</v-card-title>
      <v-container>
        <v-row>
          <v-col>
            <v-text-field
              v-model="username"
              label="Username"
              required
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-combobox
              v-model="model"
              :filter="filter"
              :hide-no-data="!search"
              :items="instruments"
              :search-input.sync="search"
              hide-selected
              multiple
              small-chips
              solo
              label="Select instrument(s) that you play"
            >
              <template v-slot:no-data>
                <v-list-item>
                  <span class="subheading">Create unlisted instrument</span>
                </v-list-item>
                <v-chip
                  label
                  small
                >
                  {{ search }}
                </v-chip>
              </template>
              <template v-slot:selection="{ attrs, item, parent, selected }">
                <v-chip
                  v-if="item === Object(item)"
                  v-bind="attrs"
                  :input-value="selected"
                  label
                  small
                >
                  <span class="pr-2">
                    {{ item.text }}
                  </span>
                  <v-icon
                    small
                    @click="parent.selectItem(item)"
                  >mdi-close</v-icon>
                </v-chip>
              </template>
            </v-combobox>
          </v-col>
        </v-row>
      </v-container>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue darken-1" text @click="emitCreateProfile()">Save Profile</v-btn>
      </v-card-actions>
    </v-card>
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import users from "@/store/modules/users";

@Component
export default class CreateProfileDialog extends Vue {
  createProfile = false;
  @Prop() address!: string
  username = "";
  editing: Record<string, any> = {};
  index = -1;
  nonce = 1;
  search = "Guitar";
  instruments: Array<object> = [
    {
      header: "Add instrument(s)"
    },
    {
      text: "Guitar",
    },
    {
      text: "Drums"
    },
    {
      text: "Bass"
    },
    {
      text: "Vocals"
    },
    {
      text: "Piano"
    }
  ];
  model: Array<Record<string, string>> = [];

  @Watch("model")
  modelChanged(val: Array<Record<string,string>>, prev: Array<Record<string, string>>) {
    console.log("model changed")
    console.log(val)
    console.log(prev)

    if (val.length === prev.length) {
      return;
    }

    this.model = val.map((v: Record<string, string>) => {
      if (typeof v === 'string') {
        v = {
          text: v
        }
        this.instruments.push(v)
        this.nonce++
      }
      return v;
    });
  }

  filter (item: {header: string, text: string}, queryText: string, itemText: string) {
    console.log("filter")
    if (item.header) return false

    const hasValue = (val: string) => val != null ? val : '';

    const text = hasValue(itemText);
    const query = hasValue(queryText);

    return text.toString()
      .toLowerCase()
      .indexOf(query.toString().toLowerCase()) > -1
  }

  handleCreateProfile() {
    console.log("create profile request");
  }
}

</script>