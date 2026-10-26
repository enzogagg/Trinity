import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: { proxy: { "/api": "http://api:8080" } },
  test: {
    coverage: {
      provider: "v8",
      include: ["src/**/*.{ts,tsx}"],
      // cobertura : lu par GitLab pour afficher la couverture dans les merge requests
      reporter: ["text", "text-summary", "cobertura"],
    },
  },
});
