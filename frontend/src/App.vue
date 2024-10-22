<script setup lang="ts">
	import { onMounted, ref } from "vue";

	import Navbar from "@/components/NavBar.vue";
	import { useWebsocket, type IWorkerConfig } from "@/stores/websocket";

	const ws = useWebsocket();
	const config = ref<IWorkerConfig>({ endpoint: "/api/classes", parameters: "warlock" });
	const resp = ref({});

	onMounted(() => {
		ws.Connect();
		FetchAPI();
	});

	function FetchAPI() {
		ws.CallAPI({ ...config.value });

		// Implicitly not awaiting it so the result is stored concurrently
		ws.ListenForResponse().then(result => resp.value = JSON.parse(result));
	}
</script>

<template>
	<header>
		<Navbar />
	</header>

	<main>
		Is connected ? {{ ws.inner.connected }}
		<button @click="FetchAPI">Fetch API</button>
		<div>{{ resp }}</div>
	</main>
</template>

<style scoped lang="scss">
</style>
