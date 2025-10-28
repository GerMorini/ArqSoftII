/**
 * Hook personalizado para manejar actividades
 * Centraliza toda la lógica de fetch y estado de actividades
 */

import { useState, useEffect, useCallback } from 'react';
import { actividadService } from '../services/actividadService';
import logger from '../utils/logger';

export function useActividades(usuarioId = null) {
  const [actividades, setActividades] = useState([]);
  const [inscripciones, setInscripciones] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  /**
   * Cargar todas las actividades
   */
  const fetchActividades = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await actividadService.getActividades();
      setActividades(data.activities || []);
    } catch (err) {
      const errorMessage = err.message || 'Error desconocido al cargar actividades';
      setError(errorMessage);
      logger.error('useActividades - fetch error', err);
    } finally {
      setLoading(false);
    }
  }, []);

  /**
   * Cargar inscripciones del usuario
   */
  const fetchInscripciones = useCallback(async (uid) => {
    if (!uid) return;
    try {
      const data = await actividadService.getInscripciones(uid);
      setInscripciones(data || []);
    } catch (err) {
      logger.error('useActividades - fetch inscripciones error', err);
      // No setear error global para no interrumpir el flujo
    }
  }, []);

  /**
   * Crear nueva actividad
   */
  const createActividad = useCallback(async (actividadData) => {
    try {
      const newActividad = await actividadService.createActividad(actividadData);
      setActividades((prev) => [...prev, newActividad]);
      return newActividad;
    } catch (err) {
      logger.error('useActividades - create error', err);
      throw err;
    }
  }, []);

  /**
   * Actualizar actividad existente
   */
  const updateActividad = useCallback(async (actividadId, actividadData) => {
    try {
      const updatedActividad = await actividadService.updateActividad(actividadId, actividadData);
      setActividades((prev) =>
        prev.map((act) => (act.id === actividadId ? updatedActividad : act))
      );
      return updatedActividad;
    } catch (err) {
      logger.error('useActividades - update error', err);
      throw err;
    }
  }, []);

  /**
   * Eliminar actividad
   */
  const deleteActividad = useCallback(async (actividadId) => {
    try {
      await actividadService.deleteActividad(actividadId);
      setActividades((prev) => prev.filter((act) => act.id !== actividadId));
      return true;
    } catch (err) {
      logger.error('useActividades - delete error', err);
      throw err;
    }
  }, []);

  /**
   * Inscribir usuario en actividad
   */
  const enrollInActividad = useCallback(async (usuarioIdParam, actividadId) => {
    try {
      await actividadService.enrollInActividad(usuarioIdParam, actividadId);
      // Recargar inscripciones del usuario
      await fetchInscripciones(usuarioIdParam);
      return true;
    } catch (err) {
      logger.error('useActividades - enroll error', err);
      throw err;
    }
  }, [fetchInscripciones]);

  /**
   * Desincribir usuario de actividad
   */
  const unenrollFromActividad = useCallback(async (usuarioIdParam, actividadId) => {
    try {
      await actividadService.unenrollFromActividad(usuarioIdParam, actividadId);
      // Recargar inscripciones del usuario
      await fetchInscripciones(usuarioIdParam);
      return true;
    } catch (err) {
      logger.error('useActividades - unenroll error', err);
      throw err;
    }
  }, [fetchInscripciones]);

  /**
   * Verificar si un usuario está inscrito en una actividad
   */
  const estaInscripto = useCallback((actividadId) => {
    return inscripciones.includes(actividadId);
  }, [inscripciones]);

  return {
    actividades,
    inscripciones,
    loading,
    error,
    fetchActividades,
    fetchInscripciones,
    createActividad,
    updateActividad,
    deleteActividad,
    enrollInActividad,
    unenrollFromActividad,
    estaInscripto,
  };
}
