import { ref } from "vue";
import { defineStore } from "pinia";
import { io as socketio } from "socket.io-client";

const [WEBSOCKET_URI, WEBSOCKET_PORT, WEBSOCKET_TIMEOUT, WEBSOCKET_RECONNECTION_TIMES] = [
	import.meta.env.VITE_WEBSOCKET_URI,
	import.meta.env.VITE_WEBSOCKET_PORT,
	import.meta.env.VITE_WEBSOCKET_TIMEOUT,
	import.meta.env.VITE_WEBSOCKET_RECONNECTION_TIMES,
];

export const useWebsocket = defineStore("websocket", () => {
	const websocket = socketio(WEBSOCKET_URI, {
		port: WEBSOCKET_PORT,
		timeout: WEBSOCKET_TIMEOUT,
		reconnection: true,
		reconnectionDelay: WEBSOCKET_TIMEOUT,
		reconnectionAttempts: WEBSOCKET_RECONNECTION_TIMES,
	});

	function handle_error(error: Error) {
		console.error(`[WS Error]: A websocket error occured: ${error.message}`);
	}

	websocket.io.on("error", handle_error);
	websocket.io.on("reconnect_error", handle_error);
	websocket.io.on("open", () => {});
	websocket.io.on("close", (reason, desc) => {});
	websocket.io.on("reconnect", attempt => {});
	websocket.io.on("packet", packet => {});

	const ws = ref(websocket);

	return { ws };
});
