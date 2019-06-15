<template>
  <div class="container">
    <div class="area-title">ラベル</div>

    <div class="area-field">
      <div v-if="editable">
        <input v-model="label" placeholder="label">
      </div>
    </div>

    <div v-if="editable" class="area-button">
      <button class="icon-button" @click="add">
        <font-awesome-icon icon="plus-square"/>
      </button>
      <button class="icon-button" @click="cancel">
        <font-awesome-icon icon="window-close"/>
      </button>
    </div>
    <div v-else class="area-button">
      <button class="icon-button" @click="edit">
        <font-awesome-icon icon="edit"/>
      </button>
    </div>

    <div class="area-labels" v-if="book && book.labels">
      <div v-for="label in book.labels" :key="label.id">
        <Label :book-id="book.id" :label="label" :deletable="editable" @update="onUpdate" class="label"/>
      </div>
    </div>
  </div>
</template>

<style scoped>
.icon-button {
  background-color: transparent;
  text-decoration: none;
  position: relative;
  vertical-align: middle;
  text-align: center;
  display: inline-block;
  border-radius: 3rem;
  transition: all ease 0.4s;
  padding: 0.5rem 1rem;
  margin: 0 0.5rem;
  cursor: pointer;
}

.container {
  width: 64rem;
  height: 100%;
  display: grid;

  grid-template-columns: 8rem 1fr auto;
  grid-template-rows: 2rem 1fr;
  grid-template-areas: "title field button" "labels labels labels";

  justify-items: left;
  align-items: center;
}

.area-title {
  grid-area: title;
}

.area-field {
  grid-area: field;
  width: 16rem;
  justify-self: right;
  padding: 0rem 0.5rem;
}

.area-field input {
  width: 100%;
}

.area-button {
  grid-area: button;

  justify-self: right;
}

.area-labels {
  grid-area: labels;
  display: flex;
  flex-flow: row wrap;
}

.label {
  margin: 0.2rem 0.2rem;
}
</style>

<script>
import axios from "axios";
import mixin from "../mixin";
import Label from "@/components/Label";

export default {
  components: {
    Label
  },
  props: ["book"],
  mixins: [mixin],
  data() {
    return {
      editable: false,
      label: ""
    };
  },
  methods: {
    add() {
      var data = {
        label: this.label
      };

      var me = this;

      axios
        .post("/api/books/" + this.book.id + "/labels", data)
        // eslint-disable-next-line
        .then(response => {
          me.cancel();
          me.$emit("update");
        })
        .catch(error => {
          switch (error.response.status) {
            case 401:
              me.login();
              break;
            default:
              // eslint-disable-next-line
              console.log("error:", error.response.status, error.response.data);
              break;
          }
        });
    },
    cancel() {
      this.label = "";
      this.editable = false;
    },
    edit() {
      this.editable = true;
    },
    onUpdate() {
      this.cancel();
      this.$emit("update");
    }
  }
};
</script>
