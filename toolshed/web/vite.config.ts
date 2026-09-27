import { defineConfig } from "vite";

// Proxies /resources and /loans to the Go API during local dev, so the
// frontend's fetch("/resources") calls just work without CORS setup.
export default defineConfig({
  server: {
    proxy: {
      "/resources": "http://localhost:8080",
      "/loans": "http://localhost:8080",
    },
  },
});
