<template>
  <div class="container">
    <div v-if="book && book.cover" class="area-cover">
      <img :src="book.cover">
    </div>
    <div v-else class="area-cover">
      <img src="@/assets/noimage-cover.png">
    </div>

    <div class="area-title-key">タイトル</div>
    <div v-if="book" class="area-title-value">{{book.title}}</div>
    <div v-else class="area-title-value">Unknown Title</div>

    <div class="area-author-key">著者</div>
    <div v-if="book" class="area-author-value">{{book.authors.join(", ")}}</div>
    <div v-else class="area-author-value">Unknown Authors</div>

    <div class="area-pubdate-key">出版日</div>
    <div v-if="book" class="area-pubdate-value">{{formatDate}}</div>
    <div v-else class="area-pubdate-value">YYYY/MM/DD</div>

    <div class="area-publisher-key">出版社</div>
    <div v-if="book" class="area-publisher-value">{{book.publisher}}</div>
    <div v-else class="area-publisher-value">Unknown Publisher</div>

    <div class="area-isbn-key">ISBN</div>
    <div v-if="book" class="area-isbn-value">{{book.isbn}}</div>
    <div v-else class="area-isbn-value">Unknown ISBN</div>

    <div v-if="book" class="area-label">
      <div v-for="label in book.labels" :key="label.id">
        <Label :book-id="book.id" :label="label" :editable="false" class="label"/>
      </div>
    </div>
    <div v-else class="area-label">Unknown Labels</div>

    <div class="area-recommend">
      <span v-if="liked">
        <button class="button-like" @click="dislike">
          <font-awesome-icon :icon="['fas', 'heart']"/>
        </button>
      </span>
      <span v-else>
        <button class="button-like" @click="like">
          <font-awesome-icon :icon="['far', 'heart']"/>
        </button>
      </span>
      <Balloon :text="likes"/>
    </div>
  </div>
</template>

<style scoped>
.container {
  width: 64rem;
  height: 100%;
  display: grid;
  grid-template-columns: 8rem 8rem 1fr;
  grid-template-rows: 2rem 2rem 2rem 2rem 2rem 2rem 1fr;
  grid-template-areas: "cover title-key title-value" "cover author-key author-value" "cover pubdate-key pubdate-value" "cover publisher-key publisher-value" "cover isbn-key isbn-value" "cover space space" "recommend label label";

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

  font-weight: bold;
  font-size: 1.2rem;

  padding: 0rem 0.5rem;
}

.area-author-key {
  grid-area: author-key;

  padding: 0rem 0.5rem;
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

.area-label {
  grid-area: label;
  display: flex;
  flex-flow: row wrap;

  padding: 0rem 0.5rem;
}

.label {
  margin: 0.2rem 0.2rem;
}

.area-recommend {
  grid-area: recommend;

  align-self: start;
  justify-self: center;
}

.button-like {
  background-color: transparent;
  text-decoration: none;
  position: relative;
  vertical-align: middle;
  text-align: center;
  display: inline-block;
  border-radius: 3rem;
  transition: all ease 0.4s;
  padding: 0.5rem 1rem;
  cursor: pointer;
}
</style>

<script>
import axios from "axios";
import moment from "moment";
import Balloon from "@/components/Balloon";
import mixin from "../mixin";
import Label from "@/components/Label";

export default {
  components: {
    Balloon,
    Label
  },
  props: ["book"],
  mixins: [mixin],
  data() {
    return {
      likes: 0,
      userID: ""
    };
  },
  updated() {
    this.likes = this.book.recommenders.length;
  },
  computed: {
    formatDate: function() {
      return moment(this.book.pubdate).format("YYYY/MM/DD");
    },
    liked: function() {
      return this.hasLiked();
    }
  },
  methods: {
    hasLiked() {
      this.context();

      if (this.book === null) {
        return false;
      }

      var me = this;
      var liked = false;
      this.book.recommenders.forEach(element => {
        if (me.userID == element.id) {
          liked = true;
        }
      });
      return liked;
    },
    context() {
      var me = this;
      axios
        .get("/api/context")
        .then(response => {
          me.userID = response.data.id;
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
    like() {
      var me = this;

      axios
        .post("/api/books/" + this.book.id + "/recommends")
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
    dislike() {
      var me = this;

      axios
        .delete("/api/books/" + this.book.id + "/recommends")
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
