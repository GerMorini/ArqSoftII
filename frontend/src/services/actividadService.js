/**
 * Servicio centralizado para operaciones de API de actividades
 * Centraliza todas las llamadas HTTP relacionadas con actividades
 */

import config from '../config/env';
import logger from '../utils/logger';

const ACTIVITIES_URL = config.ACTIVITIES_URL;

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
        logger.logApiError('/activities', response.status, 'Failed to create actividad');
        throw new Error('Error al crear la actividad');
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
        logger.logApiError(`/activities/${actividadId}`, response.status, 'Failed to update actividad');
        throw new Error('Error al actualizar la actividad');
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
        logger.logApiError(`/activities/${actividadId}`, response.status, 'Failed to delete actividad');
        throw new Error('Error al eliminar la actividad');
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
        logger.logApiError(`/activities/${actividadId}/inscribir`, response.status, 'Failed to enroll in actividad');
        throw new Error('Error al inscribirse en la actividad');
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
        logger.logApiError(`/activities/${actividadId}/desinscribir`, response.status, 'Failed to unenroll from actividad');
        throw new Error('Error al desincribirse de la actividad');
      }

      logger.info('Usuario desincrito exitosamente', { usuarioId, actividadId });
      return true;
    } catch (error) {
      logger.error(`Error al desincribirse de actividad ${actividadId}`, error);
      throw error;
    }
  },
};
