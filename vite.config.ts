import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import fs from 'fs'

function hotFilePlugin() {
  return {
    name: 'hot-file-plugin',
    configureServer() {
      const hotFile = path.resolve('public/hot')
      fs.mkdirSync(path.dirname(hotFile), { recursive: true })
      fs.writeFileSync(hotFile, 'http://localhost:5173')
      const cleanup = () => {
        try {
          if (fs.existsSync(hotFile)) fs.unlinkSync(hotFile)
        } catch {}
      }
      process.on('exit', cleanup)
      process.on('SIGINT', () => { cleanup(); process.exit(); })
      process.on('SIGTERM', () => { cleanup(); process.exit(); })
    },
  }
}

export default defineConfig({
  plugins: [react(), hotFilePlugin()],
  publicDir: false,
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './frontend'),
    },
  },
  server: {
    port: 5173,
    strictPort: true,
    origin: 'http://localhost:5173',
  },
  build: {
    outDir: 'public/build',
    manifest: 'manifest.json',
    rollupOptions: {
      input: 'frontend/main.tsx',
    },
  },
})