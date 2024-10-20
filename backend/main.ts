import process from "node:process";
import dotenv, { type DotenvPopulateInput } from "npm:dotenv@16.4.5";
import { Server, type ServerOptions } from "https://deno.land/x/socket_io@0.2.0/mod.ts";
import { connect as createRedisClient, type RedisConnectOptions } from "https://deno.land/x/redis@v0.33.0/mod.ts";
import { serve } from "https://deno.land/std@0.224.0/http/mod.ts";

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
	const io = new Server(socket_options);
	// const client = await createRedisClient(redis_options);
	// const sub_client = await createRedisClient(redis_options);

	io.on("connection", socket => {
		console.log("new", socket.id);
		socket.emit("message", "from server");

		socket.on("message", message => {
			console.log(`socket ${socket.id} message: ${message}`);
		});

		socket.on("disconnect", reason => {
			console.log(`socket ${socket.id} disconnected due to ${reason}`);
		});
	});

	const local_addr: Deno.Addr = {
		transport: "tcp",
		hostname: process.env.VITE_WEBSOCKET_URI!,
		port: parseInt(process.env.VITE_WEBSOCKET_PORT!),
	};

	// Deno.serve({ port: local_addr.port }, async (req, info) => {
	// 	console.log("req from", info.remoteAddr.hostname, info.remoteAddr.port);
	// 	const handler = io.handler();
	// 	return await handler(req, { remoteAddr: info.remoteAddr, localAddr: local_addr });
	// });
    // TODO: Fix and replace with above if possible, serve from http lib is deprecated
    await serve(io.handler(), { port: local_addr.port });
}
