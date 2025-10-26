/**
 * Validaciones para formularios de usuarios
 */

export const validateUsuarioForm = (formData, isCreateMode = false) => {
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
