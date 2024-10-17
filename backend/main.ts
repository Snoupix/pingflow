import { Server } from "npm:socket.io@4.8.0";
import type { ServerOptions } from "npm:socket.io@4.8.0";
import { createClient } from "npm:redis@4.7.0";

const socket_options: Partial<ServerOptions> = {};

if (import.meta.main) {
	main();
}

function main() {
	const socket = new Server(socket_options);
}
