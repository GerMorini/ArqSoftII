/**
 * Logger centralizado para la aplicación
 * Proporciona un único punto de entrada para todos los logs
 */

const LOG_LEVELS = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3,
  SILENT: 4
};

// Obtener nivel de log de variable de entorno o usar 'INFO' por defecto
const LOG_LEVEL_ENV = (import.meta.env.VITE_LOG_LEVEL || 'INFO').toUpperCase();
const currentLogLevel = LOG_LEVELS[LOG_LEVEL_ENV] !== undefined ? LOG_LEVELS[LOG_LEVEL_ENV] : LOG_LEVELS.INFO;

const logger = {
  /**
   * Log de información general
   */
  info: (message, data = null) => {
    if (currentLogLevel <= LOG_LEVELS.INFO) {
      console.log(`[INFO] ${message}`, data || '');
    }
  },

  /**
   * Log de advertencia
   */
  warn: (message, data = null) => {
    if (currentLogLevel <= LOG_LEVELS.WARN) {
      console.warn(`[WARN] ${message}`, data || '');
    }
  },

  /**
   * Log de error
   */
  error: (message, error = null) => {
    if (currentLogLevel <= LOG_LEVELS.ERROR) {
      console.error(`[ERROR] ${message}`, error || '');
    }
  },

  /**
   * Log de debug (solo si está habilitado)
   */
  debug: (message, data = null) => {
    if (currentLogLevel <= LOG_LEVELS.DEBUG) {
      console.debug(`[DEBUG] ${message}`, data || '');
    }
  },

  /**
   * Log específico para operaciones de actividades
   */
  logActivityFetch: (url) => {
    logger.debug(`Fetching actividades from: ${url}`);
  },

  /**
   * Log de acciones en actividades (create, update, delete, enroll)
   */
  logActivityAction: (action, actividadId, details = null) => {
    logger.info(`Activity action [${action}] on actividad ${actividadId}`, details);
  },

  /**
   * Log de inscripciones de usuario
   */
  logInscripcion: (action, usuarioId, actividadId) => {
    logger.info(`Inscripción action [${action}]: usuario ${usuarioId} - actividad ${actividadId}`);
  },

  /**
   * Log de errores de API
   */
  logApiError: (endpoint, statusCode, errorMessage) => {
    logger.error(`API Error [${statusCode}] on ${endpoint}:`, errorMessage);
  },
};

export default logger;
