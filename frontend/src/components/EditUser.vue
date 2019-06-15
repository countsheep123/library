<template>
  <div class="container">
    <div class="area-title">ユーザ一覧</div>

    <div v-if="editable" class="area-button">
      <button class="icon-button" @click="cancel">
        <font-awesome-icon icon="window-close"/>
      </button>
    </div>
    <div v-else class="area-button">
      <button class="icon-button" @click="edit">
        <font-awesome-icon icon="edit"/>
      </button>
    </div>

    <ul class="area-users">
      <li v-for="user in users" :key="user.id">
        <User :user="user" :editable="editable" @update="onUpdate"/>
      </li>
    </ul>
  </div>
</template>

<style scoped>
ul {
  list-style: none;
  padding-left: 0;
}

li {
  margin: 0.5rem 0rem;
}

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

  grid-template-columns: 8rem 1fr 1fr 0.5fr;
  grid-template-rows: 4rem 1fr;
  grid-template-areas: "title empty1 empty2 button" "users users users users";

  justify-items: left;
  align-items: center;
}

.area-title {
  grid-area: title;
}

.area-empty1 {
  grid-area: empty1;
}

.area-empty2 {
  grid-area: empty2;
}

.area-button {
  grid-area: button;

  justify-self: right;
}

.area-users {
  grid-area: users;
  width: 100%;
}
</style>

<script>
import axios from "axios";
import mixin from "../mixin";
import User from "@/components/User";

export default {
  components: {
    User
  },
  mixins: [mixin],
  data() {
    return {
      editable: false,
      users: []
    };
  },
  mounted() {
    this.getUsers();
  },
  methods: {
    getUsers() {
      var me = this;
      axios
        .get("/api/users")
        .then(response => {
          me.users = response.data.users;
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
      this.editable = false;
    },
    edit() {
      this.editable = true;
    },
    onUpdate() {
      this.getUsers();
    }
  }
};
</script>
