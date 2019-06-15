<template>
  <div class="container">
    <v-quagga :onDetected="detected" :readerSize="readerSize" :readerTypes="['ean_reader']"></v-quagga>
  </div>
</template>

<style scoped>
.container {
  padding: 2rem;
  width: 640px;
  margin: 0 auto;
  text-align: center;
}
</style>

<script>
import axios from "axios";
import moment from "moment";

export default {
  data() {
    return {
      detactable: true,
      readerSize: {
        width: 640,
        height: 480
      }
    };
  },
  methods: {
    detected(result) {
      var me = this;
      var book = null;

      var code = result.codeResult.code;
      // var format = result.codeResult.format;
      if (this.detactable) {
        if (code.startsWith("978") || code.startsWith("979")) {
          this.detactable = false;
          Promise.all([
            axios.get("https://api.openbd.jp/v1/get?isbn=" + code),
            axios.get(
              "https://www.googleapis.com/books/v1/volumes?q=isbn:" + code
            )
          ])
            .then(([response1, response2]) => {
              if (response1.status == 200) {
                if (response1.data.length > 0 && response1.data[0] != null) {
                  var authors = [];
                  response1.data[0].onix.DescriptiveDetail.Contributor.forEach(
                    element => {
                      authors.push(element.PersonName.content);
                    }
                  );

                  var book1 = {
                    isbn: code,
                    title: response1.data[0].summary.title,
                    publisher: response1.data[0].summary.publisher,
                    pubdate: moment(response1.data[0].summary.pubdate).format(
                      "YYYYMMDD"
                    ),
                    authors: authors,
                    cover: response1.data[0].summary.cover
                  };

                  book = book1;
                }
              }

              if (response2.status == 200) {
                if (response2.data.totalItems > 0) {
                  var book2 = {
                    isbn: code,
                    title: response2.data.items[0].volumeInfo.title,
                    publisher: response2.data.items[0].volumeInfo.publisher,
                    pubdate: moment(
                      response2.data.items[0].volumeInfo.publishedDate
                    ).format("YYYYMMDD"),
                    authors: response2.data.items[0].volumeInfo.authors,
                    cover:
                      response2.data.items[0].volumeInfo.imageLinks.thumbnail
                  };

                  if (book == null) {
                    book = book2;
                  }
                }
              }

              if (book !== null) {
                me.$emit("detected", book);
              } else {
                me.$emit("failed");
              }
            })
            .catch(([error1, error2]) => {
              // eslint-disable-next-line
              console.log("error1:", error1.response.status, error1.response.data);
              // eslint-disable-next-line
              console.log("error2:", error2.response.status, error2.response.data);

              me.$emit("failed");
            });
        } else {
          // eslint-disable-next-line
          console.log("978又は979から始まるバーコードをかざしてください");
        }
      }
    }
  }
};
</script>