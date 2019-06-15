<template>
  <div class="user-container">
    <div class="area-user">{{user.name}}</div>
    <div class="area-button">
      <div v-if="editable">
        <div v-if="user.is_admin">
          <button class="icon-button" @click="lock">
            <font-awesome-icon icon="lock-open"/>
          </button>
        </div>
        <div v-else>
          <button class="icon-button" @click="unlock">
            <font-awesome-icon icon="lock"/>
          </button>
        </div>
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

.user-container {
  width: 100%;
  height: 100%;
  display: grid;
  grid-template-columns: 8rem 1fr;
  grid-template-rows: 2rem;
  grid-template-areas: "user button";

  justify-items: left;
  align-items: center;
}

.area-user {
  grid-area: user;

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
  props: ["user", "editable"],
  mixins: [mixin],
  data() {
    return {};
  },
  methods: {
    lock() {
      this.update(false);
    },
    unlock() {
      this.update(true);
    },
    update(isAdmin) {
      var me = this;

      var data = {
        is_admin: isAdmin
      };

      axios
        .patch("/api/users/" + this.user.id, data)
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
