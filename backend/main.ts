import process from "node:process";
import dotenv, { type DotenvPopulateInput } from "npm:dotenv@16.4.5";
import { Server, type ServerOptions } from "npm:socket.io@4.8.0";
import { connect as createRedisClient, type RedisConnectOptions } from "https://deno.land/x/redis@v0.33.0/mod.ts";

const socket_options: Partial<ServerOptions> = {};
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
	// const socket = new Server(socket_options);
	// const client = await createRedisClient(redis_options);
	// const sub_client = await createRedisClient(redis_options);
}
