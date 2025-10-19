// Configuración de variables de entorno
// Las variables se cargan desde el archivo .env en el directorio raíz del proyecto

const config = {
  USERS_URL: import.meta.env.VITE_USERS_URL || 'http://localhost:8083',
  ACTIVITIES_URL: import.meta.env.VITE_ACTIVITIES_URL || 'http://localhost:8084'
};

// Log para debugging (solo en desarrollo)
if (import.meta.env.DEV) {
  console.log('🔧 Variables de entorno cargadas:');
  console.log('VITE_USERS_URL:', config.USERS_URL);
  console.log('VITE_ACTIVITIES_URL:', config.ACTIVITIES_URL);
}

export default config;
