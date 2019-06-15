<template>
  <div class="location-container">
    <div class="area-location">{{location.name}}</div>
    <div class="area-button">
      <div v-if="deletable">
        <button class="icon-button" @click="confirm">
          <font-awesome-icon icon="trash-alt"/>
        </button>
      </div>
    </div>

    <dialog ref="confirm">
      本当に削除しますか？
      <button class="icon-button" @click="remove">
        <font-awesome-icon icon="check-square"/>
      </button>
      <button class="icon-button" @click="closeModal">
        <font-awesome-icon icon="window-close"/>
      </button>
    </dialog>
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

.location-container {
  width: 100%;
  height: 100%;
  display: grid;
  grid-template-columns: 10rem 1fr;
  grid-template-rows: 2rem;
  grid-template-areas: "location button";

  justify-items: left;
  align-items: center;
}

.area-location {
  grid-area: location;

  padding: 0rem 0.5rem;
}

.area-button {
  grid-area: button;

  justify-self: right;
}
</style>

<script>
import axios from "axios";
import mixin from "../mixin";

export default {
  components: {},
  props: ["location", "deletable"],
  mixins: [mixin],
  data() {
    return {};
  },
  methods: {
    confirm() {
      this.$refs.confirm.showModal();
    },
    remove() {
      var me = this;

      this.closeModal();

      axios
        .delete("/api/locations/" + this.location.id)
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
    },
    closeModal() {
      this.$refs.confirm.close();
    }
  }
};
</script>
