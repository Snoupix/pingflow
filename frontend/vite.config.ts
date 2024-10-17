import path from "node:path";
import dotenv from "dotenv";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue"; // https://vitejs.dev/config/ export default

dotenv.config({ path: path.resolve(__dirname, "../.env.public") });

export default defineConfig({ plugins: [vue()] });
