<template>
  <div>
    <div>
      <button class="icon-button" @click="edit">
        <font-awesome-icon icon="edit"/>
      </button>

      <span v-if="isAdmin">
        <button class="icon-button" @click="confirm">
          <font-awesome-icon icon="trash-alt"/>
        </button>
      </span>
    </div>

    <Card class="card">
      <Book :book="book" @update="onUpdate"></Book>
    </Card>

    <Card class="card">
      <EditLabel :book="book" @update="onUpdate"></EditLabel>
    </Card>

    <Card v-if="book" class="card">
      <EditStock :book="book" @update="onUpdate"></EditStock>
    </Card>

    <dialog ref="confirm">
      本当に削除しますか？
      <button class="icon-button" @click="remove(book.id)">
        <font-awesome-icon icon="check-square"/>
      </button>
      <button class="icon-button" @click="closeModal">
        <font-awesome-icon icon="window-close"/>
      </button>
    </dialog>
  </div>
</template>

<style scoped>
.card {
  display: inline-flex;
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
import EditLabel from "@/components/EditLabel";
import EditStock from "@/components/EditStock";
import mixin from "../mixin";
import axios from "axios";

export default {
  components: {
    Card,
    Book,
    EditLabel,
    EditStock
  },
  mixins: [mixin],
  data() {
    return {
      book: null,
      id: "",
      isAdmin: false
    };
  },
  created() {
    var id = this.$route.params.id;
    this.id = id;

    if (id === null || id === undefined) {
      this.$router.push({
        name: "list-book"
      });
      return;
    }

    this.fetch(id);
    this.context();
  },
  methods: {
    fetch(id) {
      var me = this;
      axios
        .get("/api/books/" + id)
        .then(response => {
          switch (response.status) {
            case 200:
              me.book = response.data;
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
    confirm() {
      this.$refs.confirm.showModal();
    },
    closeModal() {
      this.$refs.confirm.close();
    },
    edit() {
      this.$router.push({
        name: "edit",
        params: { book: this.book }
      });
    },
    remove(id) {
      var me = this;

      axios
        .delete("/api/books/" + id)
        .then(response => {
          switch (response.status) {
            case 204:
              this.$router.push({
                name: "list-book"
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
    onUpdate() {
      this.fetch(this.id);
    },
    context() {
      var me = this;
      axios
        .get("/api/context")
        .then(response => {
          me.isAdmin = response.data.is_admin;
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
