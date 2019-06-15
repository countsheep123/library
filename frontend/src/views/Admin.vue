<template>
  <div>
    <Card class="card">
      <EditMark></EditMark>
    </Card>

    <div v-if="isAdmin">
      <Card class="card">
        <EditLocation></EditLocation>
      </Card>
    </div>

    <div v-if="isAdmin">
      <Card class="card">
        <EditUser></EditUser>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.card {
  display: inline-flex;
}
</style>

<script>
import Card from "@/components/Card";
import EditMark from "@/components/EditMark";
import EditLocation from "@/components/EditLocation";
import EditUser from "@/components/EditUser";
import mixin from "../mixin";
import axios from "axios";

export default {
  components: {
    Card,
    EditMark,
    EditLocation,
    EditUser
  },
  mixins: [mixin],
  data() {
    return {
      isAdmin: false
    };
  },
  created() {
    this.context();
  },
  watch: {
    isAdmin: function(val) {
      if (!val) {
        this.$router.push({
          name: "list-book"
        });
      }
    }
  },
  methods: {
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
