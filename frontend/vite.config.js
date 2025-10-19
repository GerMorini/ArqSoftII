// vite.config.js
import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig(({ mode }) => {
  // Cargar variables de entorno del archivo .env en el directorio raíz
  const env = loadEnv(mode, '../', '');
  
  return {
    plugins: [react()],
    envDir: '../', // Buscar en el directorio padre (raíz del proyecto)
    envPrefix: 'VITE_',
    define: {
      // Usar las variables del .env o valores por defecto
      'import.meta.env.VITE_USERS_URL': JSON.stringify(env.VITE_USERS_URL || 'http://localhost:8083'),
      'import.meta.env.VITE_ACTIVITIES_URL': JSON.stringify(env.VITE_ACTIVITIES_URL || 'http://localhost:8084'),
    },
  };
});