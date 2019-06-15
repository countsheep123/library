<template>
  <div class="container">
    <div class="area-title">場所一覧</div>

    <div v-if="editable" class="area-name">
      <input v-model="name" placeholder="name">
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

    <ul class="area-locations">
      <li v-for="location in locations" :key="location.id">
        <Location :location="location" :deletable="editable" @update="onUpdate"/>
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
  grid-template-areas: "title empty name button" "locations locations locations locations";

  justify-items: left;
  align-items: center;
}

.area-title {
  grid-area: title;
}

.area-empty {
  grid-area: empty;
}

.area-name {
  grid-area: name;
  width: 100%;
}

.area-name > input {
  width: 80%;
}

.area-button {
  grid-area: button;

  justify-self: right;
}

.area-locations {
  grid-area: locations;
  width: 100%;
}
</style>

<script>
import axios from "axios";
import mixin from "../mixin";
import Location from "@/components/Location";

export default {
  components: {
    Location
  },
  mixins: [mixin],
  data() {
    return {
      editable: false,
      locations: [],
      name: ""
    };
  },
  mounted() {
    this.getLocations();
  },
  methods: {
    getLocations() {
      var me = this;
      axios
        .get("/api/locations")
        .then(response => {
          me.locations = response.data;
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
        name: this.name
      };

      var me = this;

      axios
        .post("/api/locations", data)
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
      this.editable = false;
      this.getLocations();
    },
    edit() {
      this.editable = true;
    },
    onUpdate() {
      this.cancel();
    }
  }
};
</script>
