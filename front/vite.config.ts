import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  preview: {
    allowedHosts: ["kakeibo.taesubase.com", "localhost"],
    port: 3000,
    host: true,
  }
})
