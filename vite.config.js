import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";

export default defineConfig({
  plugins: [solidPlugin({})],
  build: {
    manifest: true,
    rollupOptions: {
      input: "standards/client/index.tsx",
      output: {
        dir: "dist/client",
        entryFileNames: "[name].js",
      },
    },
  },
});
