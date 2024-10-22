import { ref } from "vue";
import { defineStore } from "pinia";
import { io as socketio } from "socket.io-client";

import type { SocketOptions, ManagerOptions } from "socket.io-client";

const WEBSOCKET_OPTIONS: Partial<ManagerOptions & SocketOptions> = {
	autoConnect: false,
	reconnection: true,
	timeout: import.meta.env.VITE_WEBSOCKET_TIMEOUT,
	reconnectionAttempts: import.meta.env.VITE_WEBSOCKET_RECONNECTION_TIMES,
};

export type IWorkerConfig = {
	endpoint: string;
	parameters: string;
};

/**
 * Don't forget to use the .connect() method
 */
export const useWebsocket = defineStore("websocket", () => {
	const WORKER_API_KEY = "call_api";
	const websocket = ref(socketio(
		`${import.meta.env.VITE_WEBSOCKET_URI}:${import.meta.env.VITE_WEBSOCKET_PORT}`,
		WEBSOCKET_OPTIONS,
	));

	websocket.value.on("error", error => {
		console.error(`[WS Error]: A websocket error occured ${error}`);
	});

	function CallAPI(config: IWorkerConfig) {
		websocket.value.emit(WORKER_API_KEY, JSON.stringify(config));
	}

	// Returns unparsed JSON
	function ListenForResponse(): Promise<string> {
		return new Promise(resolve => {
			websocket.value.on(WORKER_API_KEY, result => {
				websocket.value.off(WORKER_API_KEY);
				return resolve(result);
			});
		});
	}

	function Connect() {
		websocket.value.connect();
	}

	return { inner: websocket, CallAPI, ListenForResponse, Connect };
});

/*
    const socket = websocket(`${process.env[ENV_VARS.WS_URI]}:${process.env[ENV_VARS.WS_PORT]}`, websocket_options);

    let response = "";

    socket.on("connect", () => {
        socket.send("from client");
        socket.emit("call_api", JSON.stringify({ endpoint: "/api/classes", parameters: "warlock" }));
    });

    socket.on("call_api", result => {
        // console.log("Result", result);
        response = result;
    });

    socket.on("error", err => {
        throw err;
    });

    socket.connect();

    await new Promise(resolve => setTimeout(resolve, 750));

    if (response.length == 0) {
        socket.disconnect();
        throw new Error("Failed to receive message from server in less than 750ms");
    }

    socket.disconnect();
*/
