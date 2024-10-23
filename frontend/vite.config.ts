import path from "node:path";
import { fileURLToPath, URL } from "node:url";
import dotenv from "dotenv";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueDevTools from "vite-plugin-vue-devtools";

dotenv.config({ path: path.resolve(__dirname, "../.env.public") });

// https://vite.dev/config/
export default defineConfig({
	css: {
		preprocessorOptions: {
			sass: {
				api: "modern",
			},
		},
	},
	plugins: [vue(), vueDevTools()],
	resolve: {
		alias: {
			"@": fileURLToPath(new URL("./src", import.meta.url)),
		},
	},
});
