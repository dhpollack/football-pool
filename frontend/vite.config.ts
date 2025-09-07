import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import basicSsl from "@vitejs/plugin-basic-ssl";

const plugins = [react(), tailwindcss()];
if (!process.env.RUNNING_WITH_PLAYWRIGHT) {
  plugins.unshift(basicSsl());
}

// https://vite.dev/config/
export default defineConfig({
  plugins: plugins,
  define: {
    global: {},
  },
});
