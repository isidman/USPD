import { defineConfig } from "vite";

// No backend to proxy to — the engine is pure client-side calculation,
// same as the original app ("no backend — all React state, designed for
// PWA conversion post-makeathon").
export default defineConfig({});
