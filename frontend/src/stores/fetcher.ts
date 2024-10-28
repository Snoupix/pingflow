import { ref, type Ref } from "vue";
import { defineStore } from "pinia";

import { useWebsocket, type IWorkerConfig } from "@/stores/websocket";

import type { Error as APIError, API_OUT } from "@/types/api";

export const useFetcher = defineStore("fetcher", () => {
	const error = ref<APIError | null>(null);
	const is_loading = ref(false);

	function Fetch<T extends API_OUT>(websocket: ReturnType<typeof useWebsocket>, state: Ref<T | null>, config: IWorkerConfig) {
		is_loading.value = true;

		websocket.CallAPI(config);

		// Implicitly not awaiting it so the result is stored concurrently
		websocket.ListenForResponse<T>()
			.then(resp => (state.value = resp))
			.catch(err => (error.value = err))
			.finally(() => (is_loading.value = false))
	}

	return { error, is_loading, Fetch };
});
