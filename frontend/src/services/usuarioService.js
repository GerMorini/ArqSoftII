/**
 * Servicio centralizado para operaciones de usuarios
 * Centraliza todas las llamadas HTTP relacionadas con usuarios y provee utilidades
 * para la manipulacion de los mismos
 */

import config from '../config/env';
import logger from '../utils/logger';
import { getTokenPayload } from '../utils/tokenUtils';

const USERS_URL = config.USERS_URL;

/**
 * Limpia todos los datos de sesión del usuario
 */
const clearUserSession = () => {
  localStorage.removeItem('access_token');
  logger.info('Sesión de usuario limpiada');
};

/**
 * Guardar sesión de usuario en localStorage
 * Solo guarda el JWT token - los datos del usuario se decodifican on-demand
 */
const storeUserSession = (accessToken, username) => {
  try {
    localStorage.setItem('access_token', accessToken);
    logger.info('Sesión de usuario guardada', { username });
  } catch (error) {
    logger.error('Error al guardar sesión de usuario', error);
    throw error;
  }
};

/**
 * Validar formato de email
 * @param {string} email - Email a validar
 * @returns {object} { valid: boolean, error?: string }
 */
const validateEmail = (email) => {
  if (!email || email.trim() === '') {
    return { valid: false, error: 'El email es requerido' };
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(email)) {
    return { valid: false, error: 'El formato del email no es válido' };
  }

  return { valid: true };
};

/**
 * Validar contraseña
 * @param {string} password - Contraseña a validar
 * @returns {object} { valid: boolean, error?: string }
 */
const validatePassword = (password) => {
  if (!password || password === '') {
    return { valid: false, error: 'La contraseña es requerida' };
  }

  if (password.length < 6) {
    return { valid: false, error: 'La contraseña debe tener al menos 6 caracteres' };
  }

  // Validar que contenga mayúscula, minúscula y número
  const strongRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{6,}$/;
  if (!strongRegex.test(password)) {
    return { valid: false, error: 'La contraseña debe contener mayúscula, minúscula y número' };
  }

  return { valid: true };
};

/**
 * Validaciones para formularios de usuarios
 */
const validateUsuarioForm = (formData, isCreateMode = false) => {
  const errors = {};

  // Nombre
  if (!formData.nombre || formData.nombre.trim() === '') {
    errors.nombre = 'El nombre es requerido';
  } else if (formData.nombre.trim().length < 2) {
    errors.nombre = 'El nombre debe tener al menos 2 caracteres';
  }

  // Apellido
  if (!formData.apellido || formData.apellido.trim() === '') {
    errors.apellido = 'El apellido es requerido';
  } else if (formData.apellido.trim().length < 2) {
    errors.apellido = 'El apellido debe tener al menos 2 caracteres';
  }

  // Email
  if (!formData.email || formData.email.trim() === '') {
    errors.email = 'El email es requerido';
  } else if (!isValidEmail(formData.email)) {
    errors.email = 'El email no es válido';
  }

  // Validaciones solo para modo crear
  if (isCreateMode) {
    // Username
    if (!formData.username || formData.username.trim() === '') {
      errors.username = 'El nombre de usuario es requerido';
    } else if (formData.username.trim().length < 3) {
      errors.username = 'El nombre de usuario debe tener al menos 3 caracteres';
    }

    // Password
    if (!formData.password || formData.password === '') {
      errors.password = 'La contraseña es requerida';
    } else if (formData.password.length < 6) {
      errors.password = 'La contraseña debe tener al menos 6 caracteres';
    }
  } else {
    // En modo edición, validar contraseña solo si se ingresa una
    if (formData.password && formData.password.trim() !== '' && formData.password.length < 6) {
      errors.password = 'La contraseña debe tener al menos 6 caracteres';
    }
  }

  return errors;
};

const isValidEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

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

  /**
   * Iniciar sesión de usuario
   */
  async login(username, password) {
    try {
      logger.info('Intentando iniciar sesión', { username });
      const response = await fetch(`${USERS_URL}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: username.trim(),
          password
        })
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error de autenticación');
      }

      const data = await response.json();
      storeUserSession(data.access_token, username);
      logger.info('Sesión iniciada exitosamente', { username });
      return data;
    } catch (error) {
      logger.error('Error al iniciar sesión', error);
      throw error;
    }
  },

  /**
   * Registrar nuevo usuario (público)
   */
  async register(userData) {
    try {
      logger.info('Intentando registrar nuevo usuario', { username: userData.username });
      const response = await fetch(`${USERS_URL}/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          nombre: userData.nombre.trim(),
          apellido: userData.apellido.trim(),
          email: userData.email.trim(),
          username: userData.username.trim(),
          password: userData.password
        })
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error al registrar usuario');
      }

      const data = await response.json();
      storeUserSession(data.access_token, userData.username);
      logger.info('Usuario registrado exitosamente', { username: userData.username });
      return data;
    } catch (error) {
      logger.error('Error al registrar usuario', error);
      throw error;
    }
  },

  /**
   * Crear usuario desde el panel admin
   */
  async createUsuarioAdmin(userData) {
    try {
      logger.info('Creando nuevo usuario desde admin', { username: userData.username });
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${USERS_URL}/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          nombre: userData.nombre.trim(),
          apellido: userData.apellido.trim(),
          email: userData.email.trim(),
          username: userData.username.trim(),
          password: userData.password,
          is_admin: userData.is_admin || false
        })
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error al crear usuario');
      }

      const data = await response.json();
      logger.info('Usuario creado exitosamente desde admin', { username: userData.username });
      return data;
    } catch (error) {
      logger.error('Error al crear usuario desde admin', error);
      throw error;
    }
  },
  validateEmail,
  validatePassword,
  validateUsuarioForm,
  clearUserSession,
};
