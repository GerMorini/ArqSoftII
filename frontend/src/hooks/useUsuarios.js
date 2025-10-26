import { useState, useCallback } from 'react';
import { usuarioService } from '../services/usuarioService';
import logger from '../utils/logger';

/**
 * Hook personalizado para gestionar usuarios (admin)
 * Proporciona métodos CRUD y manejo de estado para usuarios
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

  return {
    usuarios,
    loading,
    error,
    fetchUsuarios,
    getUsuarioById,
    updateUsuario,
    deleteUsuario,
  };
}
