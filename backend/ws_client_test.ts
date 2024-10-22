import { assertNotEquals } from "jsr:@std/assert";
import process from "node:process";
import dotenv, { type DotenvPopulateInput } from "npm:dotenv@16.4.5";
import { io as websocket } from "npm:socket.io-client@4.8.0";
import type { SocketOptions, ManagerOptions } from "npm:socket.io-client@4.8.0";

const WEBSOCKET_OPTIONS: Partial<ManagerOptions & SocketOptions> = {
	autoConnect: false,
	reconnection: true,
};

const ENV_VARS = {
	WS_URI: "VITE_WEBSOCKET_URI",
	WS_PORT: "VITE_WEBSOCKET_PORT",
	WS_TIMEOUT: "VITE_WEBSOCKET_TIMEOUT",
	WS_RECONNECTION_TIMES: "VITE_WEBSOCKET_RECONNECTION_TIMES",
};

Deno.test("socket.io client test script", { sanitizeResources: false, sanitizeOps: false }, async t => {
	await t.step("Configuring ENV based on .env* files", () => {
		for (const path of ["../.env", "../.env.public"]) {
			const err = dotenv.config({ path, processEnv: process.env as DotenvPopulateInput }).error;
			if (err != undefined) {
				throw err;
			}
		}

		for (const env_var of Object.values(ENV_VARS)) {
			assertNotEquals(process.env[env_var], undefined);
			assertNotEquals(process.env[env_var], "");
		}

		WEBSOCKET_OPTIONS.timeout = parseInt(process.env[ENV_VARS.WS_TIMEOUT]!);
		WEBSOCKET_OPTIONS.reconnectionAttempts = parseInt(process.env[ENV_VARS.WS_RECONNECTION_TIMES]!);
	});

	await t.step("Running Websocket with socket.io and awaiting response message", async () => {
		// Since Deno doesn't implements the node:http api yet, it successully reaches (connects ?) to
		// socket.io server but messages are not received/sent
		//
		// ref https://github.com/denoland/deno/issues/19507
		// ref https://docs.deno.com/runtime/reference/node_apis/#node%3Ahttp
		const socket = websocket(`${process.env[ENV_VARS.WS_URI]}:${process.env[ENV_VARS.WS_PORT]}`, WEBSOCKET_OPTIONS);

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
	});
});
