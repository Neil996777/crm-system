import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/auth': {
        target: process.env.VITE_GATEWAY_PROXY_TARGET ?? 'http://127.0.0.1:8080',
        changeOrigin: true
      },
      '/admin': {
        target: process.env.VITE_GATEWAY_PROXY_TARGET ?? 'http://127.0.0.1:8080',
        changeOrigin: true
      },
      '/api': {
        target: process.env.VITE_GATEWAY_PROXY_TARGET ?? 'http://127.0.0.1:8080',
        changeOrigin: true
      }
    }
  }
});
