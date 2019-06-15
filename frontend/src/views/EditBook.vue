<template>
  <div>
    <div v-if="enable">
      <div class="scan">
        <button class="icon-button" @click="stopScan">
          <font-awesome-icon icon="ban"/>
        </button>
        <BarcodeReader @detected="onDetected"/>
      </div>
    </div>

    <div v-else>
      <div>
        <div v-if="isUpdate">
          <button class="icon-button" @click="edit">
            <font-awesome-icon icon="check-square"/>
          </button>
          <button class="icon-button" @click="cancel">
            <font-awesome-icon icon="window-close"/>
          </button>
        </div>
        <div v-else>
          <button class="icon-button" @click="add">
            <font-awesome-icon icon="plus-square"/>
          </button>
          <button class="icon-button" @click="startScan">
            <font-awesome-icon icon="barcode"/>
          </button>
        </div>
      </div>

      <Card class="card">
        <div class="container">
          <div v-if="data && data.cover" class="area-cover">
            <img :src="data.cover">
          </div>
          <div v-else class="area-cover">
            <img src="@/assets/noimage-cover.png">
          </div>

          <label class="area-title-key" for="title">タイトル</label>
          <input
            class="area-title-value"
            type="text"
            id="title"
            v-model="data.title"
            placeholder="title"
          >

          <div class="area-author-key">
            <label for="authors">著者</label>
            <font-awesome-icon icon="plus-square" @click="addAuthor" class="icon"/>
          </div>
          <div class="area-author-value">
            <div v-for="(author, index) in data.authors" :key="index">
              <input type="text" id="authors" v-model="author.value" placeholder="authors">
              <font-awesome-icon icon="minus-square" @click="removeAuthor(index)" class="icon"/>
            </div>
          </div>

          <label class="area-pubdate-key" for="pubdate">出版日</label>
          <input
            class="area-pubdate-value"
            type="text"
            id="pubdate"
            v-model="data.pubdate"
            placeholder="pubdate"
          >

          <label class="area-publisher-key" for="publisher">出版社</label>
          <input
            class="area-publisher-value"
            type="text"
            id="publisher"
            v-model="data.publisher"
            placeholder="publisher"
          >

          <label class="area-isbn-key" for="isbn">ISBN</label>
          <input
            class="area-isbn-value"
            type="text"
            id="isbn"
            v-model="data.isbn"
            placeholder="isbn"
            :readonly="isUpdate"
          >
        </div>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.scan {
  display: flex;
  flex-flow: column nowrap;
}

.card {
  display: inline-flex;
}

.icon {
  cursor: pointer;
  margin: 0 0.5rem;
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
  grid-template-columns: 8rem 8rem 1fr;
  grid-template-rows: 2rem auto 2rem 2rem 2rem 2rem;
  grid-template-areas: "cover title-key title-value" "cover author-key author-value" "cover pubdate-key pubdate-value" "cover publisher-key publisher-value" "cover isbn-key isbn-value" "cover label label";

  justify-items: left;
  align-items: center;
}

.area-cover {
  grid-area: cover;
  padding: 0.5rem;
}

.area-cover > img {
  object-fit: cover;
  width: 100%;
  max-height: 100%;
}

.area-title-key {
  grid-area: title-key;

  padding: 0rem 0.5rem;
}

.area-title-value {
  grid-area: title-value;

  padding: 0rem 0.5rem;
}

.area-author-key {
  grid-area: author-key;

  padding: 0rem 0.5rem;
  align-items: flex-start;
}

.area-author-value {
  grid-area: author-value;

  padding: 0rem 0.5rem;
}

.area-pubdate-key {
  grid-area: pubdate-key;

  padding: 0rem 0.5rem;
}

.area-pubdate-value {
  grid-area: pubdate-value;

  padding: 0rem 0.5rem;
}

.area-publisher-key {
  grid-area: publisher-key;

  padding: 0rem 0.5rem;
}

.area-publisher-value {
  grid-area: publisher-value;

  padding: 0rem 0.5rem;
}

.area-isbn-key {
  grid-area: isbn-key;

  padding: 0rem 0.5rem;
}

.area-isbn-value {
  grid-area: isbn-value;

  padding: 0rem 0.5rem;
}
</style>

<script>
import Card from "@/components/Card";
import BarcodeReader from "@/components/BarcodeReader";
import mixin from "../mixin";
import axios from "axios";
import moment from "moment";

export default {
  components: {
    Card,
    BarcodeReader
  },
  mixins: [mixin],
  props: ["book"],
  data() {
    return {
      isUpdate: false,
      enable: false,
      data: null
    };
  },
  created() {
    this.init(this.$route.name);
  },
  methods: {
    init(name) {
      switch (name) {
        case "add":
          this.isUpdate = false;

          this.data = {
            title: "",
            publisher: "",
            pubdate: "",
            authors: [""],
            isbn: "",
            cover: ""
          };
          break;
        case "edit":
          if (this.book == null) {
            this.$router.push({
              name: "list-book"
            });
            return;
          }

          this.isUpdate = true;

          this.data = this.refresh(this.book);

          break;
        default:
          this.$router.push({
            name: "list-book"
          });
      }
    },
    onDetected(result) {
      this.enable = false;

      this.data = this.refresh(result);
    },
    startScan() {
      this.enable = true;
      this.data = {
        title: "",
        publisher: "",
        pubdate: "",
        authors: [""],
        isbn: "",
        cover: ""
      };
    },
    stopScan() {
      this.enable = false;
    },
    add() {
      var me = this;
      let book = this.convert(this.data);

      axios
        .post("/api/books", book)
        .then(response => {
          switch (response.status) {
            case 201:
              this.$router.push({
                name: "detail-book",
                params: { id: response.data.id }
              });
              break;
          }
        })
        .catch(error => {
          switch (error.response.status) {
            case 400:
              // eslint-disable-next-line
              console.log("fail:", error.response.data);
              break;
            case 401:
              me.login();
              break;
            case 409:
              // eslint-disable-next-line
              console.log("fail:", error.response.data);

              this.$router.push({
                name: "detail-book",
                params: { id: error.response.data.id }
              });

              break;
            default:
              // eslint-disable-next-line
              console.log("error:", error.response.status, error.response.data);
              break;
          }
        });
    },
    edit() {
      var me = this;
      let book = this.convert(this.data);

      axios
        .patch("/api/books/" + book.id, book)
        .then(response => {
          switch (response.status) {
            case 200:
              this.$router.push({
                name: "detail-book",
                params: { id: book.id }
              });
              break;
          }
        })
        .catch(error => {
          switch (error.response.status) {
            case 400:
              // eslint-disable-next-line
              console.log("fail:", error.response.data);
              break;
            case 401:
              me.login();
              break;
            case 409:
              // eslint-disable-next-line
              console.log("fail:", error.response.data);
              break;
            default:
              // eslint-disable-next-line
              console.log("error:", error.response.status, error.response.data);
              break;
          }
        });
    },
    cancel() {
      this.$router.push({
        name: "detail-book",
        params: { id: this.data.id }
      });
    },
    refresh(book) {
      let data = {
        id: book.id,
        title: book.title,
        publisher: book.publisher,
        pubdate: moment(book.pubdate).format("YYYY/MM/DD"),
        authors: [],
        isbn: book.isbn,
        cover: book.cover
      };

      if (book.authors != null) {
        // eslint-disable-next-line
        book.authors.forEach(function(element, index, array) {
          data.authors.push({ value: element });
        });
      }
      return data;
    },
    convert(data) {
      let book = {
        id: data.id,
        title: data.title,
        publisher: data.publisher,
        pubdate: moment(data.pubdate).format("YYYYMMDD"),
        // eslint-disable-next-line
        authors: data.authors.map(function(element, index, array) {
          return element.value;
        }),
        isbn: data.isbn,
        cover: data.cover
      };
      return book;
    },
    addAuthor() {
      this.data.authors.push({ value: "" });
    },
    removeAuthor(index) {
      this.data.authors.splice(index, 1);
    }
  }
};
</script>
