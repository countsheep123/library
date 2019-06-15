<template>
  <button class="label">
    <span>{{label.label}}</span>
    <font-awesome-icon icon="window-close" v-if="deletable" @click="remove(label.id)"/>
  </button>
</template>

<style scoped>
.label {
  background-color: #ffea00;
  border-width: 0;
  text-decoration: none;
  position: relative;
  vertical-align: middle;
  text-align: center;
  display: inline-block;
  border-radius: 3rem;
  transition: all ease 0.4s;
  padding: 0.5rem 1rem;
  margin: 0 0.5rem;
}

.label:hover {
  outline: none;
}
.label:active {
  outline: none;
}
.label:focus {
  outline: none;
}

.label svg {
  cursor: pointer;
  margin: 0 0 0 0.5rem;
}
</style>

<script>
import axios from "axios";
import mixin from "../mixin";

export default {
  props: ["bookId", "label", "deletable"],
  mixins: [mixin],
  methods: {
    remove() {
      var me = this;

      axios
        .delete("/api/books/" + this.bookId + "/labels/" + this.label.id)
        // eslint-disable-next-line
        .then(response => {
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
    }
  }
};
</script>
