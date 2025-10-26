/**
 * Servicio centralizado para operaciones de API de usuarios
 * Centraliza todas las llamadas HTTP relacionadas con usuarios (admin)
 */

import config from '../config/env';
import logger from '../utils/logger';

const USERS_URL = config.USERS_URL;

/**
 * Servicio de usuarios
 */
export const usuarioService = {
  /**
   * Obtener todos los usuarios
   */
  async getUsuarios() {
    try {
      logger.info('Obteniendo usuarios desde:', USERS_URL);
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${USERS_URL}/users`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        logger.error('Error al obtener usuarios:', response.status);
        throw new Error('Error al cargar los usuarios');
      }

      const data = await response.json();
      const usuarios = Array.isArray(data) ? data : (data.usuarios || []);
      logger.info('Usuarios cargados exitosamente', { count: usuarios.length });
      return usuarios;
    } catch (error) {
      logger.error('Error al obtener usuarios', error);
      throw error;
    }
  },

  /**
   * Obtener un usuario por ID
   */
  async getUsuarioById(usuarioId) {
    try {
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${USERS_URL}/users/${usuarioId}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        throw new Error('Error al cargar el usuario');
      }

      const data = await response.json();
      logger.info('Usuario cargado exitosamente', { usuarioId });
      return data;
    } catch (error) {
      logger.error(`Error al obtener usuario ${usuarioId}`, error);
      throw error;
    }
  },

  /**
   * Actualizar un usuario
   */
  async updateUsuario(usuarioId, usuarioData) {
    try {
      logger.info('Actualizando usuario:', usuarioId);
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${USERS_URL}/users/${usuarioId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(usuarioData),
      });

      if (!response.ok) {
        throw new Error('Error al actualizar el usuario');
      }

      const data = await response.json();
      logger.info('Usuario actualizado exitosamente', { usuarioId });
      return data;
    } catch (error) {
      logger.error(`Error al actualizar usuario ${usuarioId}`, error);
      throw error;
    }
  },

  /**
   * Eliminar un usuario
   */
  async deleteUsuario(usuarioId) {
    try {
      logger.info('Eliminando usuario:', usuarioId);
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${USERS_URL}/users/${usuarioId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        throw new Error('Error al eliminar el usuario');
      }

      logger.info('Usuario eliminado exitosamente', { usuarioId });
      return true;
    } catch (error) {
      logger.error(`Error al eliminar usuario ${usuarioId}`, error);
      throw error;
    }
  },
};
