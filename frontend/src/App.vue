<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { storeToRefs } from "pinia";

import Navbar from "@/components/NavBar.vue";
import ClassComponent from "@/components/ClassComponent.vue";
import { useFetcher } from "@/stores/fetcher";
import { useWebsocket, type IWorkerConfig } from "@/stores/websocket";

import type { Classes, Class } from "@/types/api";

const default_cfg: IWorkerConfig = {
	endpoint: "/api/classes",
	parameters: "",
} as const;

const ws = useWebsocket();
const _fetcher = useFetcher();
const { error, is_loading } = storeToRefs(_fetcher);
const { Fetch } = _fetcher;
const is_ws_connected = ref(false);
const classes = ref<Classes | null>(null);
const selected_class = ref<Class | null>(null);

watch(ws, w => (is_ws_connected.value = w.inner.connected));

watch(error, err => {
	if (err != null) {
		selected_class.value = null;
	}
});

onMounted(() => {
	ws.Connect();
	Fetch(ws, classes, default_cfg);
});

function FetchClass(config: IWorkerConfig) {
	selected_class.value = null;
	Fetch(ws, selected_class, config);
}

function RemoveError() {
	error.value = null;
}
</script>

<template>
	<header>
		<Navbar :ws_state="is_ws_connected" :color="ws.current_color" />
	</header>

	<main>
		<h2 v-if="selected_class == null && !is_loading">Choose a class to fetch</h2>

		<div v-if="error != null" class="error">
			<h3>{{ error.error || error }}</h3>
			<button @click="RemoveError()">Close</button>
		</div>

		<section v-else-if="classes != null" class="classes_wrapper">
			<button
				v-for="resp in classes.results"
				:key="resp.index"
				:class="{ selected: selected_class && resp.url == selected_class.url }"
				@click="
					FetchClass({ endpoint: resp.url.replace(resp.index, ''), parameters: resp.index })
				">
				{{ resp.name }}
			</button>
		</section>

		<!-- Loading display + force rerender child on change -->
		<!-- NOTE: Since ClassComponent has state, rerender it on loading f*cks its state -->
		<h1 v-if="is_loading">Loading...</h1>
		<ClassComponent
			v-if="error == null && selected_class != null"
			:_class="selected_class"
			:ws />
	</main>
</template>

<style scoped lang="scss">
h1,
h2 {
	width: 100%;
	text-align: center;
	padding: 5vh 0;
}

.error {
	font-size: bold;
	color: red;
	width: 100%;
	text-align: center;
	padding: 5vh 0;
}

.classes_wrapper {
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
	gap: 10px;
	width: 100%;
	text-align: center;
	padding: 5vh 5vw;

	& > button {
		border: 1px solid var(--color-border);
		border-radius: 0.5rem;
		background: transparent;
		display: flex;
		justify-content: center;
		align-items: center;
		width: 100%;
		height: 50px;
		cursor: pointer;
		color: var(--color-text);
		transition: all var(--transition-duration);

		&:hover {
			border-color: var(--color-border-hover);
			background: var(--color-background-soft);
			color: white;
		}

		&.selected {
			box-shadow: 0 0.5em 0.5em -0.4em #f8f8f8;
		}
	}
}
</style>
