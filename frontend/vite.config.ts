import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'
import { visualizer } from 'rollup-plugin-visualizer';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    tailwindcss(),
    visualizer({
      open: true, // This will automatically open the visualization in your browser
      gzipSize: true, // Shows the gzipped size, which is what you care about most
      brotliSize: true, // Shows brotli size as well (another compression algorithm)
    }),

  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },

  assetsInclude: ['**/*.lottie', '**/*.wasm'],
})
