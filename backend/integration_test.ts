import { assert, assertNotEquals } from "jsr:@std/assert";
import process from "node:process";
import dotenv, { type DotenvPopulateInput } from "npm:dotenv@16.4.5";
import { connect as createRedisClient, type RedisConnectOptions } from "https://deno.land/x/redis@v0.33.0/mod.ts";

const ENV_VARS = {
	WORKER_ADDR: "WORKER_ADDR",
	WORKER_PORT: "WORKER_PORT",
	REDIS_PORT: "REDIS_PORT",
	REDIS_ADDR: "REDIS_ADDR",
	REDIS_PASSWORD: "REDIS_PASSWORD",
	REDIS_WORK_PREFIX: "REDIS_WORK_PREFIX",
	REDIS_WORK_PROCESS: "REDIS_WORK_PROCESS",
	REDIS_WORK_RESULT: "REDIS_WORK_RESULT",
	REDIS_WORK_STATUS: "REDIS_WORK_STATUS",
	REDIS_CH_WORK_PROCESS: "REDIS_CH_WORK_PROCESS",
	REDIS_CH_WORK_RESULT: "REDIS_CH_WORK_RESULT",
};

// Resource sanitize disabled since I can't, for some reason, close/shutdown the Redis TCP connection(s)
Deno.test("Worker input to output", { sanitizeResources: false }, async t => {
	t.step("Configuring ENV based on .env* files", () => {
		for (const path of ["../.env", "../.env.public"]) {
			const err = dotenv.config({ path, processEnv: process.env as DotenvPopulateInput }).error;
			if (err != undefined) {
				throw err;
			}
		}

		for (const env_var in ENV_VARS) {
			assertNotEquals(process.env[env_var], undefined);
			assertNotEquals(process.env[env_var], "");
		}
	});

	const redis_options: RedisConnectOptions = {
		hostname: process.env[ENV_VARS.REDIS_ADDR]!,
		port: process.env[ENV_VARS.REDIS_PORT],
		password: process.env[ENV_VARS.REDIS_PASSWORD],
	};

	const client = await createRedisClient(redis_options);
	const sub_client = await createRedisClient(redis_options);

	await client.connect();
	await sub_client.connect();

	const resp = await fetch(
		`http://${process.env[ENV_VARS.WORKER_ADDR]}:${process.env[ENV_VARS.WORKER_PORT]}/job-index`,
		{
			method: "GET",
		},
	);
	const job_idx = await resp.text();

	assertNotEquals(job_idx.length, 0, "/job-index response shouldn't be empty");

	const config_key = `${process.env[ENV_VARS.REDIS_WORK_PREFIX]}:${job_idx}:${process.env[ENV_VARS.REDIS_WORK_PROCESS]}`;
	await client.hset(config_key, {
		endpoint: "/api/classes",
		parameters: "warlock",
	});

	await client.publish(process.env[ENV_VARS.REDIS_CH_WORK_PROCESS]!, job_idx);

	const sub = await sub_client.subscribe(process.env[ENV_VARS.REDIS_CH_WORK_RESULT]!);

	let result_key = "";
	for await (const { message } of sub.receive()) {
		result_key = message;
		break;
	}

	sub.close();

	const result = await client.get(result_key);

	assertNotEquals(result, null);
	assertNotEquals(result!.trim().length, 0);
	assertNotEquals(result, "nil");

	client.close();
	sub_client.close();

	assert(client.isClosed);
	assert(sub_client.isClosed);

	// console.log(result);
});
