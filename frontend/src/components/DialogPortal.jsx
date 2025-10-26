import { createPortal } from 'react-dom';

/**
 * Portal para renderizar diálogos, alertas y modales fuera del árbol de React normal
 * Esto asegura que se muestren correctamente sin ser afectados por estilos CSS del contenedor padre
 */
const DialogPortal = ({ children }) => {
  return createPortal(children, document.body);
};

export default DialogPortal;
