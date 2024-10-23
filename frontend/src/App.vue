<script setup lang="ts">
import { onMounted, ref, watch } from "vue";

import Navbar from "@/components/NavBar.vue";
import ClassComponent from "@/components/ClassComponent.vue";
import { useWebsocket, type IWorkerConfig } from "@/stores/websocket";

import type { Classes, Class } from "./types/api";

const default_cfg: IWorkerConfig = { endpoint: "/api/classes", parameters: "" } as const;

const ws = useWebsocket();
const is_ws_connected = ref(false);
const is_loading = ref(false);
const config = ref<IWorkerConfig>(default_cfg);
const api_resp = ref<Classes | null>(null);
const api_err = ref<unknown | null>(null);
const selected_class = ref<Class | null>(null);

watch(ws, w => (is_ws_connected.value = w.inner.connected));

onMounted(() => {
	ws.Connect();
	FetchAPI(default_cfg);
});

function FetchAPI(cfg: IWorkerConfig) {
	config.value = cfg;

	ws.CallAPI(cfg);

	// Implicitly not awaiting it so the result is stored concurrently
	if (default_cfg.endpoint == cfg.endpoint && default_cfg.parameters == cfg.parameters) {
		ws.ListenForResponse<Classes>()
			.then(resp => (api_resp.value = resp))
			.catch(error => (api_err.value = error));
	} else {
		selected_class.value = null;
		is_loading.value = true;

		ws.ListenForResponse<Class>()
			.then(resp => (selected_class.value = resp))
			.catch(error => (api_err.value = error))
			.finally(() => (is_loading.value = false));
	}
}
</script>

<template>
	<header>
		<Navbar :ws_state="is_ws_connected" />
	</header>

	<main>
		<h2>Choose a class to fetch</h2>
		<input v-model="config.endpoint" placeholder="endpoint" type="text" />
		<input v-model="config.parameters" placeholder="parameters" type="text" />
		<button @click="FetchAPI(default_cfg)">Fetch API</button>
		<div v-if="api_err != null">
			<h3 class="error">{{ api_err }}</h3>
		</div>

		<section class="classes_wrapper" v-if="api_err == null && api_resp != null">
			<button
				v-for="resp in api_resp.results"
				:key="resp.index"
				:class="{ selected: selected_class && resp.url == selected_class.url }"
				@click="FetchAPI({ endpoint: resp.url.replace(resp.index, ''), parameters: resp.index })">
				{{ resp.name }}
			</button>
		</section>

		<h1 v-if="is_loading">Loading...</h1>
		<ClassComponent v-else-if="selected_class != null" :_class="selected_class" />
	</main>
</template>

<style scoped lang="scss">
.error {
	font-size: bold;
	color: red;
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
		border-radius: .5rem;
		background: transparent;
		display: flex;
		justify-content: center;
		align-items: center;
		width: 100%;
		height: 10vh;
		cursor: pointer;
		color: var(--color-text);
		transition: all 0.3s;

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
