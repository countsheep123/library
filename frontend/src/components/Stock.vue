<template>
  <div class="stock-container">
    <div class="area-name">{{stock.owner_name}}</div>
    <div v-if="stock && stock.mark_url" class="area-icon">
      <img class="mark" :src="stock.mark_url">
    </div>
    <div v-else class="area-icon">
      <img class="mark" src="@/assets/noimage-mark.png">
    </div>
    <div class="area-mark">{{stock.mark_name}}</div>
    <div class="area-location">{{stock.location_name}}</div>
    <div class="area-button">
      <div v-if="deletable">
        <button class="icon-button" @click="confirm">
          <font-awesome-icon icon="trash-alt"/>
        </button>
      </div>
      <div v-else>
        <button v-if="available" class="icon-button" @click="lent" title="借りる">
          借りる
          <font-awesome-icon icon="file-export"/>
        </button>
        <button v-else class="icon-button" @click="returned" title="返す">
          返す
          <font-awesome-icon icon="file-import"/>
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

.mark {
  width: 2rem;
  height: 2rem;
}

.stock-container {
  width: 100%;
  height: 100%;
  display: grid;
  grid-template-columns: 8rem 2rem 8rem 10rem 1fr;
  grid-template-rows: 2rem;
  grid-template-areas: "name icon mark location button";

  justify-items: left;
  align-items: center;
}

.area-name {
  grid-area: name;

  padding: 0rem 0.5rem;
}

.area-icon {
  grid-area: icon;
  max-height: 100%;
}

.area-icon > img {
  object-fit: cover;
  width: 100%;
  max-height: 100%;
  margin: 0.1rem;
}

.area-mark {
  grid-area: mark;

  padding: 0rem 0.5rem;
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
  props: ["stock", "bookId", "available", "recordId", "deletable"],
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
        .delete("/api/books/" + this.bookId + "/stocks/" + this.stock.id)
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
    },
    lent() {
      var me = this;

      var data = {
        stock_id: this.stock.id
      };

      axios
        .post("/api/records", data)
        .then(response => {
          me.recordId = response.data.id;

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
    returned() {
      var me = this;

      axios
        .delete("/api/records/" + this.recordId)
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
