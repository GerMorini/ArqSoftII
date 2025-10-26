/**
 * Hook personalizado para manejo de formularios con validación
 * Centraliza la lógica de estado y validación de formularios
 */

import { useState } from 'react';

export function useFormWithValidation(initialValues, onSubmit, validateFn) {
  const [formData, setFormData] = useState(initialValues);
  const [validationErrors, setValidationErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState(null);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    // Limpiar error del campo cuando el usuario empieza a escribir
    if (validationErrors[name]) {
      setValidationErrors((prev) => ({
        ...prev,
        [name]: undefined,
      }));
    }
  };

  const validateForm = () => {
    const errors = validateFn(formData);
    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitError(null);

    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);
    try {
      await onSubmit(formData);
    } catch (error) {
      setSubmitError(error.message || 'Error al procesar el formulario');
    } finally {
      setIsSubmitting(false);
    }
  };

  const resetForm = () => {
    setFormData(initialValues);
    setValidationErrors({});
    setSubmitError(null);
  };

  const updateFormData = (updates) => {
    setFormData((prev) => ({
      ...prev,
      ...updates,
    }));
  };

  return {
    formData,
    validationErrors,
    isSubmitting,
    submitError,
    handleChange,
    handleSubmit,
    resetForm,
    updateFormData,
    setFormData,
    setValidationErrors,
  };
}
