import { ref } from "vue";
import { defineStore } from "pinia";
import { io as socketio } from "socket.io-client";

import type { SocketOptions, ManagerOptions } from "socket.io-client";

import type { Error as APIError, API_OUT } from "@/types/api";

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

export type IColor = {
	r: number;
	g: number;
	b: number;
};

/**
 * Don't forget to use the .connect() method
 */
export const useWebsocket = defineStore("websocket", () => {
	const WORKER_API_KEY = "call_api";
	const websocket = ref(
		socketio(`${import.meta.env.VITE_WEBSOCKET_URI}:${import.meta.env.VITE_WEBSOCKET_PORT}`, WEBSOCKET_OPTIONS),
	);
	const current_color = ref<IColor | null>(null);

	websocket.value.on("error", error => {
		console.error(`[WS Error]: A websocket error occured ${error}`);
	});

	websocket.value.on("color", color => {
		// Shouldn't ever throw
		current_color.value = JSON.parse(color);
	});

	function CallAPI(config: IWorkerConfig) {
		websocket.value.emit(WORKER_API_KEY, JSON.stringify(config));
	}

	function ListenForResponse<T extends API_OUT>(): Promise<T> {
		return new Promise((resolve, reject) => {
			websocket.value.on(WORKER_API_KEY, result => {
				websocket.value.off(WORKER_API_KEY);
				try {
					const json = JSON.parse(result);
					if (json.error) {
						if (json.error == "timeout") { // TODO: Impl retry
							return reject(json as APIError);
						}

						console.error(json);
						return reject(json as APIError);
					}

					return resolve(json as T);
				} catch (error: unknown) {
					return reject(error);
				}
			});
		});
	}

	function Connect() {
		websocket.value.connect();
	}

	return { inner: websocket, current_color, CallAPI, ListenForResponse, Connect };
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
