// vite.config.js
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  envDir: '../', // Ruta relativa al vite.config.js
  envPrefix: 'VITE_',
});