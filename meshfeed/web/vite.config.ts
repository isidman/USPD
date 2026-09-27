import { defineConfig } from "vite";

// Proxies /feeds to the Go node's debug API during local dev.
export default defineConfig({
  server: {
    proxy: {
      "/feeds": "http://localhost:8080",
    },
  },
});
