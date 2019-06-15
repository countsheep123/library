<template>
  <div class="container">
    <div class="area-title">マーク一覧</div>

    <div v-if="editable" class="area-name">
      <input v-model="name" placeholder="name">
    </div>
    <div v-if="editable" class="area-url">
      <input type="file" @change="fileChanged">
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

    <ul class="area-marks">
      <li v-for="mark in marks" :key="mark.id">
        <Mark :mark="mark" :deletable="editable" @update="onUpdate"/>
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

.mark {
  width: 2rem;
  height: 2rem;
}

.container {
  width: 64rem;
  height: 100%;
  display: grid;

  grid-template-columns: 8rem 1fr 1fr 0.5fr;
  grid-template-rows: 4rem 1fr;
  grid-template-areas: "title name url button" "marks marks marks marks";

  justify-items: left;
  align-items: center;
}

.area-title {
  grid-area: title;
}

.area-name {
  grid-area: name;
  width: 100%;
}

.area-name > input {
  width: 80%;
}

.area-url {
  grid-area: url;
  width: 100%;
}

.area-url > input {
  width: 80%;
}

.area-button {
  grid-area: button;

  justify-self: right;
}

.area-marks {
  grid-area: marks;
  width: 100%;
}
</style>

<script>
import axios from "axios";
import mixin from "../mixin";
import Mark from "@/components/Mark";

export default {
  components: {
    Mark
  },
  mixins: [mixin],
  data() {
    return {
      editable: false,
      marks: [],
      name: "",
      url: ""
    };
  },
  mounted() {
    this.getMarks();
  },
  methods: {
    getMarks() {
      var me = this;
      axios
        .get("/api/marks?own=true")
        .then(response => {
          me.marks = response.data;
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
    add() {
      var data = {
        name: this.name,
        url: this.url
      };

      var me = this;

      axios
        .post("/api/marks", data)
        // eslint-disable-next-line
        .then(response => {
          me.cancel();
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
      this.name = "";
      this.url = "";
      this.editable = false;
      this.getMarks();
    },
    edit() {
      this.editable = true;
    },
    onUpdate() {
      this.cancel();
    },
    fileChanged(e) {
      e.preventDefault();

      var file = e.target.files[0];
      if (!file) {
        return;
      }

      this.convert(file);
    },
    convert(file) {
      var me = this;

      var reader = new FileReader();
      reader.onload = function(e) {
        me.url = e.target.result;
      };
      reader.readAsDataURL(file);
    }
  }
};
</script>
