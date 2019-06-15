<template>
  <span class="balloon">{{text}}</span>
</template>

<style scoped>
.balloon {
  position: relative;
  background: #ffffff;
  border: 2px solid #000000;
  border-radius: 4px;
  margin: 0 0.5rem;
  padding: 0 1rem;
}

.balloon:after,
.balloon:before {
  right: 100%;
  top: 50%;
  border: solid transparent;
  content: " ";
  height: 0;
  width: 0;
  position: absolute;
  pointer-events: none;
}

.balloon:after {
  border-right-color: #000000;
  border-width: 4px;
  margin-top: -4px;
}
.balloon:before {
  border-right-color: #000000;
  border-width: 7px;
  margin-top: -7px;
}
</style>

<script>
export default {
  props: ["text"],
};
</script>