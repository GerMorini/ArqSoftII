import { useState, useCallback } from 'react';
import { usuarioService } from '../services/usuarioService';
import logger from '../utils/logger';

/**
 * Hook personalizado para gestionar usuarios
 * Proporciona métodos CRUD, autenticación y manejo de estado para usuarios
 */
export function useUsuarios() {
  const [usuarios, setUsuarios] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchUsuarios = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await usuarioService.getUsuarios();
      setUsuarios(data);
      return data;
    } catch (err) {
      logger.error('Error en fetchUsuarios:', err);
      setError(err.message || 'Error al cargar usuarios');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const getUsuarioById = useCallback(async (usuarioId) => {
    try {
      const data = await usuarioService.getUsuarioById(usuarioId);
      return data;
    } catch (err) {
      logger.error('Error en getUsuarioById:', err);
      throw err;
    }
  }, []);

  const updateUsuario = useCallback(async (usuarioId, usuarioData) => {
    setLoading(true);
    setError(null);
    try {
      const data = await usuarioService.updateUsuario(usuarioId, usuarioData);
      // Actualizar el estado local
      setUsuarios(prev =>
        prev.map(u => u.id_usuario === usuarioId ? data : u)
      );
      return data;
    } catch (err) {
      logger.error('Error en updateUsuario:', err);
      setError(err.message || 'Error al actualizar usuario');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const deleteUsuario = useCallback(async (usuarioId) => {
    setLoading(true);
    setError(null);
    try {
      await usuarioService.deleteUsuario(usuarioId);
      // Actualizar el estado local removiendo el usuario eliminado
      setUsuarios(prev => prev.filter(u => u.id_usuario !== usuarioId));
      return true;
    } catch (err) {
      logger.error('Error en deleteUsuario:', err);
      setError(err.message || 'Error al eliminar usuario');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  /**
   * Iniciar sesión
   */
  const login = useCallback(async (username, password) => {
    setLoading(true);
    setError(null);
    try {
      const data = await usuarioService.login(username, password);
      return data;
    } catch (err) {
      logger.error('Error en login:', err);
      setError(err.message || 'Error al iniciar sesión');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  /**
   * Registrar nuevo usuario
   */
  const register = useCallback(async (userData) => {
    setLoading(true);
    setError(null);
    try {
      const data = await usuarioService.register(userData);
      return data;
    } catch (err) {
      logger.error('Error en register:', err);
      setError(err.message || 'Error al registrar usuario');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  /**
   * Cerrar sesión
   */
  const logout = useCallback(() => {
    try {
      localStorage.removeItem('access_token');
      localStorage.removeItem('idUsuario');
      localStorage.removeItem('isAdmin');
      localStorage.removeItem('isLoggedIn');
      localStorage.removeItem('username');
      logger.info('Sesión cerrada');
    } catch (err) {
      logger.error('Error en logout:', err);
      throw err;
    }
  }, []);

  /**
   * Obtener usuario actual desde localStorage
   */
  const getCurrentUser = useCallback(() => {
    const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';
    if (!isLoggedIn) {
      return null;
    }
    return {
      id_usuario: parseInt(localStorage.getItem('idUsuario')),
      username: localStorage.getItem('username'),
      is_admin: localStorage.getItem('isAdmin') === 'true'
    };
  }, []);

  return {
    usuarios,
    loading,
    error,
    fetchUsuarios,
    getUsuarioById,
    updateUsuario,
    deleteUsuario,
    login,
    register,
    logout,
    getCurrentUser,
  };
}
