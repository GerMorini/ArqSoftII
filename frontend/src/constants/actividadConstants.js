/**
 * Constantes para el módulo de actividades
 * Centraliza todas las constantes reutilizadas
 */

export const DIAS_SEMANA = [
  { value: 'Lunes', label: 'Lunes' },
  { value: 'Martes', label: 'Martes' },
  { value: 'Miércoles', label: 'Miércoles' },
  { value: 'Jueves', label: 'Jueves' },
  { value: 'Viernes', label: 'Viernes' },
  { value: 'Sábado', label: 'Sábado' },
  { value: 'Domingo', label: 'Domingo' },
];

export const ACTIVIDAD_STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive',
  FULL: 'full',
};

export const ACTIVIDAD_ACTIONS = {
  CREATE: 'create',
  UPDATE: 'update',
  DELETE: 'delete',
  ENROLL: 'enroll',
  UNENROLL: 'unenroll',
};

export const ERROR_MESSAGES = {
  FETCH_ACTIVIDADES: 'Error al cargar las actividades',
  FETCH_INSCRIPCIONES: 'Error al cargar las inscripciones',
  CREATE_ACTIVIDAD: 'Error al crear la actividad',
  UPDATE_ACTIVIDAD: 'Error al actualizar la actividad',
  DELETE_ACTIVIDAD: 'Error al eliminar la actividad',
  ENROLL_ACTIVIDAD: 'Error al inscribirse en la actividad',
  UNENROLL_ACTIVIDAD: 'Error al desincribirse de la actividad',
  VALIDATION_ERROR: 'Por favor revisa los campos del formulario',
  NETWORK_ERROR: 'Error de conexión. Intenta más tarde',
};

export const SUCCESS_MESSAGES = {
  ACTIVIDAD_CREATED: 'Actividad creada exitosamente',
  ACTIVIDAD_UPDATED: 'Actividad actualizada exitosamente',
  ACTIVIDAD_DELETED: 'Actividad eliminada exitosamente',
  INSCRIPCION_SUCCESS: 'Te has inscrito exitosamente',
  DESINCRIPCION_SUCCESS: 'Te has desincrito exitosamente',
};
