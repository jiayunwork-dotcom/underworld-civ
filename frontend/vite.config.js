import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte()],
  server: {
    port: 3000,
    host: '0.0.0.0',
    proxy: {
      '/api': 'http://backend:8080',
      '/ws': {
        target: 'ws://backend:8080',
        ws: true
      }
    }
  }
})
