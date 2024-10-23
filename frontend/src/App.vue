<script setup lang="ts">
import { onMounted, ref, watch } from "vue";

import Navbar from "@/components/NavBar.vue";
import ClassComponent from "@/components/ClassComponent.vue";
import { useWebsocket, type IWorkerConfig } from "@/stores/websocket";

import type { Error as APIError, Classes, Class, Spell, API_OUT_T, SpellResp } from "./types/api";

const default_cfg: IWorkerConfig & { _type: API_OUT_T } = {
	endpoint: "/api/classes",
	parameters: "",
	_type: "classes",
} as const;

const ws = useWebsocket();
const is_ws_connected = ref(false);
const is_loading = ref(false);
const config = ref<IWorkerConfig>(default_cfg);
const api_resp = ref<Classes | null>(null);
const api_err = ref<APIError | null>(null);
const selected_class = ref<Class | null>(null);
const class_spells = ref<Array<Spell> | null>(null);

watch(ws, w => (is_ws_connected.value = w.inner.connected));

watch(api_err, err => {
	if (err != null) {
		selected_class.value = null;
		class_spells.value = null;
	}
});

watch(ws.current_color, console.log);

onMounted(() => {
	ws.Connect();
	FetchAPI(default_cfg);
});

function FetchAPI(cfg: IWorkerConfig & { _type: API_OUT_T }) {
	config.value = cfg;

	ws.CallAPI(cfg);

	// Implicitly not awaiting it so the result is stored concurrently
	switch (cfg._type) {
		case "classes":
			ws.ListenForResponse<Classes>()
				.then(resp => (api_resp.value = resp))
				.catch(error => (api_err.value = error));
			break;
		case "class":
			selected_class.value = null;
			class_spells.value = null;
			is_loading.value = true;

			ws.ListenForResponse<Class>()
				.then(resp => (selected_class.value = resp))
				.catch(error => (api_err.value = error))
				.finally(() => (is_loading.value = false));
			break;
		case "spells":
			class_spells.value = null;
			is_loading.value = true;

			ws.ListenForResponse<SpellResp>()
				.then(resp => (class_spells.value = resp.results))
				.catch(error => (api_err.value = error))
				.finally(() => (is_loading.value = false));
			break;
		default:
			console.error("Unexpected error: API output _type not handled", cfg);
	}
}

function RemoveError() {
	api_err.value = null;
}
</script>

<template>
	<header>
		<Navbar :ws_state="is_ws_connected" :color="ws.current_color" />
	</header>

	<main>
		<h2 v-if="selected_class == null">Choose a class to fetch</h2>

		<div v-if="api_err != null" class="error">
			<h3>{{ api_err.error || api_err }}</h3>
			<button @click="RemoveError()">Close</button>
		</div>

		<section class="classes_wrapper" v-else-if="api_resp != null">
			<button
				v-for="resp in api_resp.results"
				:key="resp.index"
				:class="{ selected: selected_class && resp.url == selected_class.url }"
				@click="
					FetchAPI({ endpoint: resp.url.replace(resp.index, ''), parameters: resp.index, _type: 'class' })
				">
				{{ resp.name }}
			</button>
		</section>

		<!-- Loading display + force rerender child on change -->
		<h1 v-if="is_loading">Loading...</h1>
		<ClassComponent
			v-else-if="api_err == null && selected_class != null"
			:_class="selected_class"
			:FetchAPI="FetchAPI"
			:spells="class_spells" />
	</main>
</template>

<style scoped lang="scss">
h1, h2 {
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
		height: 10vh;
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
