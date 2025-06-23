import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react-swc";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
	plugins: [tailwindcss(), react()],
	server: {
		host: true,
		proxy: {
			// Proxying requests from /api to your Go API
			"/api": {
				target: "http://localhost:1323", // 'api' is the service name, 1323 is its port
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, ""), // If your Go API doesn't have /api prefix
			},
		},
	},
});
