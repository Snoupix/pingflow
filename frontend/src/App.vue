<script setup lang="ts">
	import { onMounted, ref, watch } from "vue";

	import Navbar from "@/components/NavBar.vue";
	import { useWebsocket, type IWorkerConfig } from "@/stores/websocket";

	const ws = useWebsocket();
	const config = ref<IWorkerConfig>({ endpoint: "/api/classes", parameters: "" });
	const api_response = ref({});
	const is_ws_connected = ref(false);

	watch(ws, w => is_ws_connected.value = w.inner.connected);

	onMounted(() => {
		ws.Connect();
		FetchAPI();
	});

	function FetchAPI() {
		ws.CallAPI({ ...config.value });

		// Implicitly not awaiting it so the result is stored concurrently
		ws.ListenForResponse().then(result => api_response.value = JSON.parse(result));
	}
</script>

<template>
	<header>
		<Navbar :ws_state="is_ws_connected" />
	</header>

	<main>
		Is connected ? {{ ws.inner.connected }}
		<button @click="FetchAPI">Fetch API</button>
		<div>{{ api_response }}</div>
	</main>
</template>

<style scoped lang="scss">
</style>
