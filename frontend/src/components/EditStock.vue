<template>
  <div class="container">
    <div class="area-title">所有者一覧 ({{book.stocks.length}})</div>

    <div class="area-mark">
      <div v-if="editable">
        <v-select :options="marks" label="name" v-model="mark" placeholder="mark">
          <template slot="option" slot-scope="option">
            <div class="option">
              <span v-if="option && option.url">
                <img :src="option.url" class="mark">
              </span>
              <span v-else>
                <img src="@/assets/noimage-mark.png" class="mark">
              </span>
              <span class="option-text">{{ option.name }}</span>
            </div>
          </template>
        </v-select>
      </div>
    </div>
    <div class="area-location">
      <div v-if="editable">
        <v-select :options="locations" label="name" v-model="location" placeholder="location"></v-select>
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

    <ul class="area-owners">
      <li v-for="stock in book.stocks" :key="stock.id">
        <Stock
          :stock="stock"
          :book-id="book.id"
          :available="stock.is_available"
          :record-id="hasRecord(stock.id)"
          :deletable="editable"
          @update="onUpdate"
        />
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
  grid-template-areas: "title mark location button" "owners owners owners owners";

  justify-items: left;
  align-items: center;
}

.area-title {
  grid-area: title;
}

.area-mark {
  grid-area: mark;
  width: 100%;
}

.area-location {
  grid-area: location;
  width: 100%;
}

.area-button {
  grid-area: button;

  justify-self: right;
}

.area-owners {
  grid-area: owners;
  width: 100%;
}

.option {
  display: flex;
  align-items: center;
}

.option-text {
  margin-left: 0.5rem;
}
</style>

<script>
import axios from "axios";
import mixin from "../mixin";
import Stock from "@/components/Stock";

export default {
  components: {
    Stock
  },
  props: ["book"],
  mixins: [mixin],
  data() {
    return {
      editable: false,
      marks: [],
      locations: [],
      mark: null,
      location: null,
      stockId: "",
      records: []
    };
  },
  mounted() {
    this.getMarks();
    this.getLocations();

    this.getRecords();
  },
  methods: {
    getRecords() {
      var me = this;
      axios
        .get("/api/records?own=true&borrowed=true&book_id=" + this.book.id)
        .then(response => {
          me.records = response.data.records;
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
    hasRecord(stockID) {
      var recordID = "";

      // eslint-disable-next-line
      this.records.forEach(function(element, index, array) {
        if (stockID === element.stock_id) {
          recordID = element.id;
        }
      });

      return recordID;
    },
    add() {
      var data = {
        mark_id: this.mark.id,
        location_id: this.location.id
      };

      var me = this;

      axios
        .post("/api/books/" + this.book.id + "/stocks", data)
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
      this.mark = null;
      this.location = null;
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
