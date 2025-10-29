import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { usuarioService } from '../services/usuarioService';
import logger from '../utils/logger';

/**
 * Componente que protege rutas requiriendo permisos de administrador
 * Valida el estado de admin del usuario contra el servidor
 */
const ProtectedRoute = ({ children }) => {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);

  useEffect(() => {
    const validateAdminAccess = async () => {
      try {
        setIsLoading(true);
        const isAdmin = await usuarioService.checkAdminStatus();

        if (isAdmin) {
          setIsAuthorized(true);
        } else {
          logger.warn('Acceso denegado: usuario no es administrador');
          navigate('/login');
        }
      } catch (error) {
        logger.error('Error validando acceso de admin', error);
        navigate('/login');
      } finally {
        setIsLoading(false);
      }
    };

    validateAdminAccess();
  }, [navigate]);

  // Mostrar carga mientras se valida
  if (isLoading) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh'
      }}>
        <p>Validando permisos...</p>
      </div>
    );
  }

  // Solo renderizar si está autorizado
  return isAuthorized ? children : null;
};

export default ProtectedRoute;
