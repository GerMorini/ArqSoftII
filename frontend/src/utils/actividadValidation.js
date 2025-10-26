/**
 * Utilidades de validación para formularios de actividades
 * Centraliza toda la lógica de validación de actividades
 */

export const validateActividadForm = (formData) => {
  const errors = {};

  // Validar título
  if (!formData.titulo || !formData.titulo.trim()) {
    errors.titulo = 'El título es requerido';
  } else if (formData.titulo.trim().length < 3) {
    errors.titulo = 'El título debe tener al menos 3 caracteres';
  }

  // Validar descripción
  if (!formData.descripcion || !formData.descripcion.trim()) {
    errors.descripcion = 'La descripción es requerida';
  } else if (formData.descripcion.trim().length < 10) {
    errors.descripcion = 'La descripción debe tener al menos 10 caracteres';
  }

  // Validar instructor
  if (!formData.instructor || !formData.instructor.trim()) {
    errors.instructor = 'El instructor es requerido';
  }

  // Validar hora de inicio
  if (!formData.hora_inicio) {
    errors.hora_inicio = 'La hora de inicio es requerida';
  }

  // Validar hora de fin
  if (!formData.hora_fin) {
    errors.hora_fin = 'La hora de fin es requerida';
  }

  // Validar que hora_fin > hora_inicio
  if (formData.hora_inicio && formData.hora_fin) {
    if (formData.hora_fin <= formData.hora_inicio) {
      errors.hora_fin = 'La hora de fin debe ser posterior a la hora de inicio';
    }
  }

  // Validar día
  if (!formData.dia) {
    errors.dia = 'Debe seleccionar un día';
  }

  // Validar cupo
  if (!formData.cupo) {
    errors.cupo = 'El cupo es requerido';
  } else if (isNaN(formData.cupo) || parseInt(formData.cupo) < 1) {
    errors.cupo = 'El cupo debe ser un número mayor a 0';
  }

  // Validar foto_url (opcional, pero si se proporciona debe ser válida)
  if (formData.foto_url && typeof formData.foto_url === 'string') {
    if (formData.foto_url.trim() && !isValidUrl(formData.foto_url)) {
      errors.foto_url = 'Debe proporcionar una URL válida';
    }
  }

  return errors;
};

/**
 * Valida si una URL es válida
 */
const isValidUrl = (string) => {
  try {
    new URL(string);
    return true;
  } catch (_) {
    return false;
  }
};

/**
 * Verifica si hay errores de validación
 */
export const hasValidationErrors = (errors) => {
  return Object.keys(errors).length > 0;
};
