/**
 * Hook personalizado para manejar la tecla Escape
 * Reutilizable en múltiples componentes
 * @param {function} onEscape - Callback a ejecutar cuando se presiona Escape
 */

import { useEffect } from 'react';

export function useEscapeKey(onEscape) {
  useEffect(() => {
    const handleEscapeKey = (event) => {
      if (event.key === 'Escape') {
        onEscape();
      }
    };

    document.addEventListener('keydown', handleEscapeKey);

    return () => {
      document.removeEventListener('keydown', handleEscapeKey);
    };
  }, [onEscape]);
}
