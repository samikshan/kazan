<template>
  <v-container fluid>
    <v-combobox
      v-model="model"
      :filter="filter"
      :hide-no-data="!search"
      :items="instruments"
      :search-input.sync="search"
      hide-selected
      label="Select/add instruments in recorded track"
      multiple
      small-chips
      solo
    >
      <template v-slot:no-data>
        <v-list-item>
          <span class="subheading">Create unlisted instrument</span>
          <v-chip
            label
            small
          >
            {{ search }}
          </v-chip>
        </v-list-item>
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
      <template v-slot:item="{ index, item }">
        <v-text-field
          v-if="editing === item"
          v-model="editing.text"
          autofocus
          flat
          background-color="transparent"
          hide-details
          solo
          @keyup.enter="edit(index, item)"
        ></v-text-field>
        <v-chip
          v-else
          label
          small
        >
          {{ item.text }}
        </v-chip>
        <v-spacer></v-spacer>
        <v-list-item-action @click.stop>
          <v-btn
            icon
            @click.stop.prevent="edit(index, item)"
          >
            <v-icon>{{ editing !== item ? 'mdi-pencil' : 'mdi-check' }}</v-icon>
          </v-btn>
        </v-list-item-action>
      </template>
    </v-combobox>
  </v-container>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch, Emit } from "vue-property-decorator";

@Component
export default class TrackInfoInput extends Vue {
  editing: Record<string, any> = {};
  index = -1;
  nonce = 1;
  search = "";
  instruments: Array<object> = [
    {
      header: "Select a tag or create one"
    },
    {
      text: "Guitar",
    },
    {
      text: "Drums"
    },
    {
      text: "Vocals"
    },
    {
      text: "Bass"
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

    this.$emit("tag-added", val[val.length - 1].text);

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

  edit(index: number, item: object) {
    console.log("edit")
    if (!this.editing) {
      this.editing = item;
      this.index = index;
    } else {
      this.editing = {};
      this.index = -1
    }
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
}
</script>