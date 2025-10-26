import { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';
import '../styles/PageTransition.css';

/**
 * Componente que proporciona un efecto de fade-in automático cuando cambia de página
 * Envuelve el contenido de las páginas y aplica la animación al montar/cambiar de ruta
 */
const PageTransition = ({ children }) => {
  const location = useLocation();
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    // Reiniciar la animación cuando cambia la ubicación
    setIsVisible(false);

    // Pequeño delay para asegurar que el navegador re-pinta antes de aplicar la animación
    const timer = setTimeout(() => {
      setIsVisible(true);
    }, 10);

    return () => clearTimeout(timer);
  }, [location.pathname]);

  return (
    <div className={`page-transition ${isVisible ? 'visible' : ''}`}>
      {children}
    </div>
  );
};

export default PageTransition;
