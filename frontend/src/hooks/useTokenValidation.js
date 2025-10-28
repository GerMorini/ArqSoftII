import { useState, useEffect } from 'react';
import { isTokenExpired } from '../utils/tokenUtils';
import { usuarioService } from '../services/usuarioService';
import useCurrentUser from './useCurrentUser';
import logger from '../utils/logger';

/**
 * Hook para validar la expiración del token en cada carga de página
 * Muestra un diálogo si el token está expirado
 * La redirección al login ocurre cuando el usuario hace click en "Aceptar"
 * @param {function} setAlertDialog - Función para mostrar el diálogo de alerta
 */
export function useTokenValidation(setAlertDialog) {
  const [hasChecked, setHasChecked] = useState(false);
  const { isLoggedIn } = useCurrentUser();

  useEffect(() => {
    // Solo verificar una vez por carga de página
    if (hasChecked) {
      return;
    }

    const token = localStorage.getItem('access_token');

    // Si no hay token o no está logueado, no hay nada que verificar
    if (!token || !isLoggedIn) {
      setHasChecked(true);
      return;
    }

    // Verificar si el token está expirado
    if (isTokenExpired(token)) {
      logger.warn('Token expirado, limpiando sesión');

      // Limpiar la sesión inmediatamente
      usuarioService.clearUserSession();

      // Mostrar diálogo de expiración
      // El redirect ocurrirá cuando el usuario haga click en "Aceptar"
      setAlertDialog({
        title: 'Sesión Expirada',
        message: 'Tu sesión ha expirado. Por favor, inicia sesión nuevamente.',
        type: 'error',
        isTokenExpired: true
      });
    }

    setHasChecked(true);
  }, [setAlertDialog, hasChecked, isLoggedIn]);
}
