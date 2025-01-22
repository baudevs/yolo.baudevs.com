import { defineConfig } from 'vite'
import path from 'path'

export default defineConfig({
  root: '.',
  base: '/',
  build: {
    outDir: '../internal/web/embed',
    assetsDir: 'assets',
    emptyOutDir: true,
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      d3: path.resolve(__dirname, 'node_modules/d3/dist/d3.min.js'),
      THREE: path.resolve(__dirname, 'node_modules/three/build/three.min.js'),
    }
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:4010',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://localhost:4010',
        ws: true,
      }
    }
  }
})
