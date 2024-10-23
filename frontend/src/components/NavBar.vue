<script setup lang="ts">
import { ref, watch } from "vue";

import type { IColor } from "@/stores/websocket";

const animation_str = "move-gradient 1.75s cubic-bezier(.79,.14,.15,.86) infinite";

const props = defineProps<{
	ws_state: boolean;
	color: IColor | null;
}>();

const previous_color = ref<IColor | null>(null);

watch(
	() => props.color,
	(_, prev) => (previous_color.value = prev),
);

function ParseColor(c: IColor): string {
	return `rgb(${c.r}, ${c.g}, ${c.b})`;
}

function GetGradient(c: IColor): string {
	const prev = previous_color.value == null ? "var(--color-text)" : ParseColor(previous_color.value);
	const curr = ParseColor(c);

	return `linear-gradient(to right in hsl shorter hue, ${prev}, ${curr})`;
}
</script>

<template>
	<nav>
		<h2
			:style="{
				background: props.color != null ? GetGradient(props.color) : 'var(--color-text)',
				animation: props.color != null ? animation_str : '',
			}">
			Dungeons & Dragons Classes viewer
		</h2>
		<div class="status_wrapper">
			<span>Status: <span v-if="props.ws_state">Connected</span><span v-else>Disconnected</span></span>
			<div class="ws_status" :class="{ active: props.ws_state }"></div>
		</div>
	</nav>
</template>

<style scoped lang="scss">
h2 {
	color: transparent;
	background-size: 200%;
	background-clip: text !important;
}

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
