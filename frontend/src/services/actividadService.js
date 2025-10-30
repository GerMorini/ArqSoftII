/**
 * Servicio centralizado para operaciones de API de actividades
 * Centraliza todas las llamadas HTTP relacionadas con actividades
 */

import config from '../config/env';
import logger from '../utils/logger';

const ACTIVITIES_URL = config.ACTIVITIES_URL;

/**
 * Valida si una URL es válida
 */
const isValidUrl = (string) => {
  try {
    new URL(string);
    return true;
  } catch (_) {
    return false;
  }
};

/**
 * Valida el formulario de actividades
 */
const validateActividadForm = (formData) => {
  const errors = {};

  // Validar título
  if (!formData.titulo || !formData.titulo.trim()) {
    errors.titulo = 'El título es requerido';
  } else if (formData.titulo.trim().length < 3) {
    errors.titulo = 'El título debe tener al menos 3 caracteres';
  }

  // Validar descripción
  if (!formData.descripcion || !formData.descripcion.trim()) {
    errors.descripcion = 'La descripción es requerida';
  } else if (formData.descripcion.trim().length < 10) {
    errors.descripcion = 'La descripción debe tener al menos 10 caracteres';
  }

  // Validar instructor
  if (!formData.instructor || !formData.instructor.trim()) {
    errors.instructor = 'El instructor es requerido';
  }

  // Validar hora de inicio
  if (!formData.hora_inicio) {
    errors.hora_inicio = 'La hora de inicio es requerida';
  }

  // Validar hora de fin
  if (!formData.hora_fin) {
    errors.hora_fin = 'La hora de fin es requerida';
  }

  // Validar que hora_fin > hora_inicio
  if (formData.hora_inicio && formData.hora_fin) {
    if (formData.hora_fin <= formData.hora_inicio) {
      errors.hora_fin = 'La hora de fin debe ser posterior a la hora de inicio';
    }
  }

  // Validar día
  if (!formData.dia) {
    errors.dia = 'Debe seleccionar un día';
  }

  // Validar cupo
  if (!formData.cupo) {
    errors.cupo = 'El cupo es requerido';
  } else if (isNaN(formData.cupo) || parseInt(formData.cupo) < 1) {
    errors.cupo = 'El cupo debe ser un número mayor a 0';
  }

  // Validar foto_url (opcional, pero si se proporciona debe ser válida)
  if (formData.foto_url && typeof formData.foto_url === 'string') {
    if (formData.foto_url.trim() && !isValidUrl(formData.foto_url)) {
      errors.foto_url = 'Debe proporcionar una URL válida';
    }
  }

  return errors;
};

/**
 * Verifica si hay errores de validación
 */
const hasValidationErrors = (errors) => {
  return Object.keys(errors).length > 0;
};

/**
 * Extrae el mensaje de error desde la respuesta del backend
 * @param {Response} response - La respuesta HTTP
 * @param {string} defaultMessage - Mensaje por defecto si no se puede extraer
 * @returns {Promise<string>} - El mensaje de error
 */
async function extractErrorMessage(response, defaultMessage) {
  try {
    const data = await response.json();

    return data.error + ": " + data.details;
  } catch {
    return defaultMessage;
  }
}

/**
 * Servicio de actividades
 */
export const actividadService = {
  /**
   * Obtener todas las actividades
   */
  async getActividades() {
    try {
      logger.logActivityFetch(ACTIVITIES_URL);
      const response = await fetch(`${ACTIVITIES_URL}/activities`);

      if (!response.ok) {
        logger.logApiError('/activities', response.status, 'Failed to fetch actividades');
        throw new Error('Error al cargar las actividades');
      }

      const data = await response.json();
      logger.info('Actividades cargadas exitosamente', { count: data.activities?.length || 0 });
      return data;
    } catch (error) {
      logger.error('Error al obtener actividades', error);
      throw error;
    }
  },

  /**
   * Obtener una actividad por ID (admin view devuelve usuarios inscritos si el token es de admin)
   */
  async getActividadById(actividadId) {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${ACTIVITIES_URL}/activities/${actividadId}`, {
        method: 'GET',
        headers: token ? { 'Authorization': `Bearer ${token}` } : {}
      });

      if (!response.ok) {
        logger.logApiError(`/activities/${actividadId}`, response.status, 'Failed to fetch actividad by id');
        throw new Error('Error al cargar la actividad');
      }

      const data = await response.json();
      // API returns { activity: ... }
      return data.activity || data;
    } catch (error) {
      logger.error(`Error al obtener actividad ${actividadId}`, error);
      throw error;
    }
  },

  /**
   * Obtener inscripciones de un usuario
   */
  async getInscripciones(usuarioId) {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${ACTIVITIES_URL}/inscriptions/${usuarioId}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        logger.logApiError(`/inscriptions/${usuarioId}`, response.status, 'Failed to fetch inscripciones');
        throw new Error('Error al cargar las inscripciones');
      }

      const data = await response.json();
      const inscripciones = Array.isArray(data) ? data : (data.inscripciones || []);
      logger.info(`Inscripciones cargadas para usuario ${usuarioId}`, { count: inscripciones.length || 0 });
      return inscripciones;
    } catch (error) {
      logger.error(`Error al obtener inscripciones del usuario ${usuarioId}`, error);
      throw error;
    }
  },

  /**
   * Crear una nueva actividad
   */
  async createActividad(actividadData) {
    try {
      logger.logActivityAction('CREATE', 'new', actividadData);
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${ACTIVITIES_URL}/activities`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(actividadData),
      });

      if (!response.ok) {
        const errorMessage = await extractErrorMessage(response, 'Error al crear la actividad');
        logger.logApiError('/activities', response.status, errorMessage);
        throw new Error(errorMessage);
      }

      const data = await response.json();
      logger.info('Actividad creada exitosamente', { actividadId: data.id });
      return data;
    } catch (error) {
      logger.error('Error al crear actividad', error);
      throw error;
    }
  },

  /**
   * Actualizar una actividad existente
   */
  async updateActividad(actividadId, actividadData) {
    try {
      logger.logActivityAction('UPDATE', actividadId, actividadData);
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${ACTIVITIES_URL}/activities/${actividadId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(actividadData),
      });

      if (!response.ok) {
        const errorMessage = await extractErrorMessage(response, 'Error al actualizar la actividad');
        logger.logApiError(`/activities/${actividadId}`, response.status, errorMessage);
        throw new Error(errorMessage);
      }

      const data = await response.json();
      logger.info('Actividad actualizada exitosamente', { actividadId });
      return data;
    } catch (error) {
      logger.error(`Error al actualizar actividad ${actividadId}`, error);
      throw error;
    }
  },

  /**
   * Eliminar una actividad
   */
  async deleteActividad(actividadId) {
    try {
      logger.logActivityAction('DELETE', actividadId);
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${ACTIVITIES_URL}/activities/${actividadId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        const errorMessage = await extractErrorMessage(response, 'Error al eliminar la actividad');
        logger.logApiError(`/activities/${actividadId}`, response.status, errorMessage);
        throw new Error(errorMessage);
      }

      logger.info('Actividad eliminada exitosamente', { actividadId });
      return true;
    } catch (error) {
      logger.error(`Error al eliminar actividad ${actividadId}`, error);
      throw error;
    }
  },

  /**
   * Inscribir un usuario en una actividad
   */
  async enrollInActividad(usuarioId, actividadId) {
    try {
      logger.logInscripcion('ENROLL', usuarioId, actividadId);
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${ACTIVITIES_URL}/activities/${actividadId}/inscribir`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        const errorMessage = await extractErrorMessage(response, 'Error al inscribirse en la actividad');
        logger.logApiError(`/activities/${actividadId}/inscribir`, response.status, errorMessage);
        throw new Error(errorMessage);
      }

      const data = await response.json();
      logger.info('Usuario inscrito exitosamente', { usuarioId, actividadId });
      return data;
    } catch (error) {
      logger.error(`Error al inscribirse en actividad ${actividadId}`, error);
      throw error;
    }
  },

  /**
   * Desincribir un usuario de una actividad
   */
  async unenrollFromActividad(usuarioId, actividadId) {
    try {
      logger.logInscripcion('UNENROLL', usuarioId, actividadId);
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${ACTIVITIES_URL}/activities/${actividadId}/desinscribir`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        const errorMessage = await extractErrorMessage(response, 'Error al desincribirse de la actividad');
        logger.logApiError(`/activities/${actividadId}/desinscribir`, response.status, errorMessage);
        throw new Error(errorMessage);
      }

      logger.info('Usuario desincrito exitosamente', { usuarioId, actividadId });
      return true;
    } catch (error) {
      logger.error(`Error al desincribirse de actividad ${actividadId}`, error);
      throw error;
    }
  },

  /**
   * Validar formulario de actividades
   */
  validateActividadForm,

  /**
   * Verificar si hay errores de validación
   */
  hasValidationErrors,
};
