// vite.config.ts
import { defineConfig } from "file:///home/veto/ai/slape/SLaMO_Frontend/node_modules/.deno/vite@5.4.11/node_modules/vite/dist/node/index.js";
import react from "file:///home/veto/ai/slape/SLaMO_Frontend/node_modules/.deno/@vitejs+plugin-react@4.3.3/node_modules/@vitejs/plugin-react/dist/index.mjs";
import deno from "file:///home/veto/ai/slape/SLaMO_Frontend/node_modules/.deno/@deno+vite-plugin@1.0.0/node_modules/@deno/vite-plugin/dist/index.js";
import "file:///home/veto/ai/slape/SLaMO_Frontend/node_modules/.deno/react@18.3.1/node_modules/react/index.js";
import "file:///home/veto/ai/slape/SLaMO_Frontend/node_modules/.deno/react-dom@18.3.1/node_modules/react-dom/index.js";
var vite_config_default = defineConfig({
  root: "./client",
  server: {
    port: 3e3
  },
  plugins: [
    react(),
    deno()
  ],
  optimizeDeps: {
    include: ["react/jsx-runtime"]
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvaG9tZS92ZXRvL2FpL3NsYXBlL1NMYU1PX0Zyb250ZW5kXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCIvaG9tZS92ZXRvL2FpL3NsYXBlL1NMYU1PX0Zyb250ZW5kL3ZpdGUuY29uZmlnLnRzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9ob21lL3ZldG8vYWkvc2xhcGUvU0xhTU9fRnJvbnRlbmQvdml0ZS5jb25maWcudHNcIjtpbXBvcnQgeyBkZWZpbmVDb25maWcgfSBmcm9tIFwidml0ZVwiO1xuaW1wb3J0IHJlYWN0IGZyb20gXCJAdml0ZWpzL3BsdWdpbi1yZWFjdFwiO1xuaW1wb3J0IGRlbm8gZnJvbSBcIkBkZW5vL3ZpdGUtcGx1Z2luXCI7XG5cbmltcG9ydCBcInJlYWN0XCI7XG5pbXBvcnQgXCJyZWFjdC1kb21cIjtcblxuZXhwb3J0IGRlZmF1bHQgZGVmaW5lQ29uZmlnKHtcbiAgcm9vdDogXCIuL2NsaWVudFwiLFxuICBzZXJ2ZXI6IHtcbiAgICBwb3J0OiAzMDAwLFxuICB9LFxuICBwbHVnaW5zOiBbXG4gICAgcmVhY3QoKSxcbiAgICBkZW5vKCksXG4gIF0sXG4gIG9wdGltaXplRGVwczoge1xuICAgIGluY2x1ZGU6IFtcInJlYWN0L2pzeC1ydW50aW1lXCJdLFxuICB9LFxufSk7XG4iXSwKICAibWFwcGluZ3MiOiAiO0FBQXdSLFNBQVMsb0JBQW9CO0FBQ3JULE9BQU8sV0FBVztBQUNsQixPQUFPLFVBQVU7QUFFakIsT0FBTztBQUNQLE9BQU87QUFFUCxJQUFPLHNCQUFRLGFBQWE7QUFBQSxFQUMxQixNQUFNO0FBQUEsRUFDTixRQUFRO0FBQUEsSUFDTixNQUFNO0FBQUEsRUFDUjtBQUFBLEVBQ0EsU0FBUztBQUFBLElBQ1AsTUFBTTtBQUFBLElBQ04sS0FBSztBQUFBLEVBQ1A7QUFBQSxFQUNBLGNBQWM7QUFBQSxJQUNaLFNBQVMsQ0FBQyxtQkFBbUI7QUFBQSxFQUMvQjtBQUNGLENBQUM7IiwKICAibmFtZXMiOiBbXQp9Cg==
