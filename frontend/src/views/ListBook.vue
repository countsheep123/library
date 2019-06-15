<template>
  <div>
    <div v-if="enable">
      <div class="scan">
        <button class="icon-button" @click="stopScan">
          <font-awesome-icon icon="ban"/>
        </button>
        <BarcodeReader @detected="onDetected" @failed="onFailed"/>
      </div>
    </div>

    <div v-else>
      <div>
        <input type="radio" id="mine" value="true" v-model="own">
        <label for="mine">所有</label>

        <input type="radio" id="all" value="false" v-model="own">
        <label for="all">全て</label>
      </div>

      <input type="text" v-model="keyword">
      <button class="icon-button" @click="search">
        <font-awesome-icon icon="search"/>
      </button>
      <button class="icon-button" @click="startScan">
        <font-awesome-icon icon="barcode"/>
      </button>

      <div>
        <input type="radio" id="title" value="title" v-model="type">
        <label for="title">Title</label>

        <input type="radio" id="publisher" value="publisher" v-model="type">
        <label for="publisher">Publisher</label>

        <input type="radio" id="authors" value="authors" v-model="type">
        <label for="authors">Authors</label>

        <input type="radio" id="isbn" value="isbn" v-model="type">
        <label for="isbn">ISBN</label>

        <input type="radio" id="label" value="label" v-model="type">
        <label for="label">Label</label>
      </div>

      <div class="book">
        <div v-for="book in books" :key="book.id">
          <Card class="card">
            <Book @click.native="detail(book)" :book="book" @update="onUpdate"></Book>
          </Card>
        </div>
      </div>

      <div>
        <button class="icon-button" @click="previous" :disabled="hasPrevious">
          <font-awesome-icon icon="chevron-circle-left"/>
        </button>
        <span>({{start}}-{{last}} / {{total}})</span>
        <button class="icon-button" @click="next" :disabled="hasNext">
          <font-awesome-icon icon="chevron-circle-right"/>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  cursor: pointer;
  display: inline-flex;
}
.card:hover {
}
.card:active {
}
.scan {
  display: flex;
  flex-flow: column nowrap;
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
</style>

<script>
import Card from "@/components/Card";
import Book from "@/components/Book";
import BarcodeReader from "@/components/BarcodeReader";
import mixin from "../mixin";
import axios from "axios";

export default {
  components: {
    Card,
    Book,
    BarcodeReader
  },
  mixins: [mixin],
  data() {
    return {
      enable: false,
      keyword: "",
      type: "title",
      own: "true",
      books: [],
      start: 0,
      last: 0,
      total: 0,
      isbn: "",
      limit: 100,
      offset: 0
    };
  },
  computed: {
    hasPrevious: function() {
      return !this.checkPrevious();
    },
    hasNext: function() {
      return !this.checkNext();
    }
  },
  mounted() {
    this.search();
  },
  methods: {
    search() {
      var me = this;

      var query =
        "?own=" +
        this.own +
        "&sort=created_at:desc&limit=" +
        this.limit +
        "&offset=" +
        this.offset;
      if (this.keyword.length > 0) {
        query = query + "&" + this.type + "=" + this.keyword;
      }

      axios
        .get("/api/books" + query)
        .then(response => {
          switch (response.status) {
            case 200:
              me.books = [];
              me.total = response.data.total;
              if (me.total > 0) {
                me.start = me.offset + 1;
              } else {
                me.start = 0;
              }
              if (me.total > me.offset + me.limit) {
                me.last = me.offset + me.limit;
              } else {
                me.last = me.total;
              }

              response.data.books.forEach(element => {
                me.books.push(element);
              });
              break;
          }
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
    onDetected(result) {
      this.enable = false;

      if (result !== null) {
        this.isbn = result.isbn;
        this.keyword = this.isbn;
        this.type = "isbn";
      }
    },
    onFailed() {
      this.enable = false;
    },
    startScan() {
      this.enable = true;
    },
    stopScan() {
      this.enable = false;
    },
    detail(book) {
      this.$router.push({
        name: "detail-book",
        params: { id: book.id }
      });
    },
    onUpdate() {
      this.search();
    },
    previous() {
      this.offset -= this.limit;
      this.search();
    },
    next() {
      this.offset += this.limit;
      this.search();
    },
    checkPrevious() {
      if (this.offset === 0) {
        return false;
      }
      return true;
    },
    checkNext() {
      if (this.offset + this.limit >= this.total) {
        return false;
      }
      return true;
    }
  }
};
</script>
