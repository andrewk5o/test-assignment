import path from "node:path";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: { "@": path.resolve(__dirname, "./src") },
  },
  server: {
    proxy: {
      // Dev-only: forward API calls to the Go backend so no CORS is needed.
      "/api": process.env.VITE_API_PROXY ?? "http://localhost:8080",
    },
  },
});
