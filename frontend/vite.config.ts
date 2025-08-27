import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { tanstackRouter } from '@tanstack/router-plugin/vite'
import { visualizer } from 'rollup-plugin-visualizer';
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
		tailwindcss(),
		visualizer({
			open: true, // This will automatically open the visualization in your browser
			gzipSize: true, // Shows the gzipped size, which is what you care about most
			brotliSize: true, // Shows brotli size as well (another compression algorithm)
		}),

	],

	assetsInclude: ['**/*.lottie'],
	// base: "/web/",

})
