<template>
	<div class="card">
		<slot></slot>
	</div>
</template>

<style scoped>
.card {
	border-radius: 2px;
	box-shadow: 0px 2px 1px -1px rgba(0, 0, 0, 0.2),
		0px 1px 1px 0px rgba(0, 0, 0, 0.14), 0px 1px 3px 0px rgba(0, 0, 0, 0.12);
	margin: 1rem 0rem;
	padding: 1rem;
}
</style>

<script>
export default {
};
</script>
