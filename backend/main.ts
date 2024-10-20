import process from "node:process";
import dotenv, { type DotenvPopulateInput } from "npm:dotenv@16.4.5";
import { Server, type ServerOptions } from "https://deno.land/x/socket_io@0.2.0/mod.ts";
import { connect as createRedisClient, type RedisConnectOptions } from "https://deno.land/x/redis@v0.33.0/mod.ts";
import { serve } from "https://deno.land/std@0.224.0/http/mod.ts";
import z from "https://deno.land/x/zod@v3.23.8/mod.ts";

const IWorkerConfigData = z.object({
	endpoint: z.string(),
	parameters: z.string(),
});

const socket_options: Partial<ServerOptions> = {
	connectTimeout: 30_000,
};
let redis_options: RedisConnectOptions;

if (import.meta.main) {
	for (const path of ["../.env", "../.env.public"]) {
		const err = dotenv.config({ path, processEnv: process.env as DotenvPopulateInput }).error;
		if (err != undefined) {
			throw err;
		}
	}

	redis_options = {
		hostname: process.env.REDIS_ADDR!,
		port: process.env.REDIS_PORT,
		password: process.env.REDIS_PASSWORD,
	};

	// Yes, Deno has top-level awaits !
	await main();
}

async function main() {
	const socketio = new Server(socket_options);
	const client = await createRedisClient(redis_options);
	const sub_client = await createRedisClient(redis_options);

	await client.connect();
	await sub_client.connect();

	socketio.on("connection", socket => {
		// Receive from Client
		socket.on("call_api", async params => {
			let config_data: z.infer<typeof IWorkerConfigData> | null = null;
			try {
				config_data = IWorkerConfigData.parse(JSON.parse(params));
			} catch (error: unknown) {
				socket.emit("error", `Failed to parse JSON ${error}`);
				return;
			}

			const resp = await fetch(`http://${process.env.WORKER_ADDR}:${process.env.WORKER_PORT}/job-index`, {
				method: "GET",
			});
			const job_idx = await resp.text();

			const config_key = `${process.env.REDIS_WORK_PREFIX}:${job_idx}:${process.env.REDIS_WORK_PROCESS}`;
			await client.hset(config_key, config_data);

			// Trigger Worker
			await client.publish(process.env.REDIS_CH_WORK_PROCESS!, job_idx);
		});
	});

	// Receive from Worker
	(async () => {
		const chan = await sub_client.subscribe(process.env.REDIS_CH_WORK_RESULT!);

		for await (const { message: result_key } of chan.receive()) {
			const result = await client.get(result_key);
			if (!result || result.length == 0) {
				break;
			}

			socketio.emit("call_api", result);
		}

		await chan.unsubscribe(process.env.REDIS_CH_WORK_RESULT!);
		chan.close();
	})();

	console.log("Press Ctrl-C or send SIGINT to gracefully close the server");

	Deno.addSignalListener("SIGINT", () => {
		socketio.close();
		client.close();
		sub_client.close();
		Deno.exit(0);
	});

	const local_addr: Deno.Addr = {
		transport: "tcp",
		hostname: process.env.VITE_WEBSOCKET_URI!.replaceAll(/wss*:\/\//g, ""),
		port: parseInt(process.env.VITE_WEBSOCKET_PORT!),
	};

	// Deno.serve({ port: local_addr.port }, async (req, info) => {
	// 	console.log("req from", info.remoteAddr.hostname, info.remoteAddr.port);
	// 	const handler = io.handler();
	// 	return await handler(req, { remoteAddr: info.remoteAddr, localAddr: local_addr });
	// });
	// TODO: Fix and replace with above if possible, serve from http stdlib is deprecated
	await serve(socketio.handler(), { hostname: local_addr.hostname, port: local_addr.port });
}
