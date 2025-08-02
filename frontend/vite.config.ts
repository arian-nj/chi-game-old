import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { tanstackRouter } from '@tanstack/router-plugin/vite'

// https://vite.dev/config/
export default defineConfig({
	server: {
		proxy: {
			'/api': 'http://localhost:8383'
		}
	},
	plugins: [
		tanstackRouter({
			target: 'react',
			autoCodeSplitting: true,
		}),
		react(),
		tailwindcss()],

	assetsInclude: ['**/*.lottie'], // 👈 This line tells Vite not to parse .lottie files
	// base: "/web/",

})
