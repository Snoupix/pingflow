<script setup lang="ts">
const { ws_state } = defineProps<{
	ws_state: boolean;
}>();
</script>

<template>
	<nav>
		<h2>Dungeons & Dragons Classes viewer</h2>
		<div class="status_wrapper">
			<span>Status: <span v-if="ws_state">Connected</span><span v-else>Disconnected</span></span>
			<div class="ws_status" :class="{ active: ws_state }"></div>
		</div>
	</nav>
</template>

<style lang="scss">
nav {
	height: 10vh;
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: space-between;
	padding: 0 2.5vw;
	background: var(--color-background-soft);

	.status_wrapper {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 10px;
		height: 100%;

		.ws_status {
			border-radius: 50%;
			width: 10px;
			height: 10px;
			background: red;
			transition: background 0.75s;

			&.active {
				background: #1dc41d;
			}
		}
	}
}
</style>
