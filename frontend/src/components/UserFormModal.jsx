import React, { useState, useEffect } from 'react';
import '../styles/ActivityFormModal.css';
import { useEscapeKey } from '../hooks/useEscapeKey';
import { usuarioService } from '../services/usuarioService';
import logger from '../utils/logger';

const UserFormModal = ({ mode = 'create', usuario = null, onClose, onSave }) => {
    const [formData, setFormData] = useState({
        nombre: '',
        apellido: '',
        username: '',
        email: '',
        password: '',
        confirm_password: '',
        is_admin: false
    });
    const [submitError, setSubmitError] = useState('');
    const [validationErrors, setValidationErrors] = useState({});
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEscapeKey(onClose);

    // Inicializar formulario según modo
    useEffect(() => {
        if (mode === 'edit' && usuario) {
            const usuarioData = {
                nombre: usuario.nombre || '',
                apellido: usuario.apellido || '',
                username: usuario.username || '',
                email: usuario.email || '',
                password: '',
                is_admin: usuario.is_admin || false
            };
            setFormData(usuarioData);
        } else {
            setFormData({
                nombre: '',
                apellido: '',
                username: '',
                email: '',
                password: '',
                is_admin: false
            });
        }
    }, [mode, usuario]);

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));
        // Limpiar error de validación cuando el usuario modifica el campo
        if (validationErrors[name]) {
            setValidationErrors(prev => ({
                ...prev,
                [name]: undefined
            }));
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSubmitError('');

        // Validar formulario
        const isCreateMode = mode === 'create';
        const errors = usuarioService.validateUsuarioForm(formData, isCreateMode);
        setValidationErrors(errors);

        if (Object.keys(errors).length > 0) {
            return;
        }

        setIsSubmitting(true);

        try {
            const token = localStorage.getItem('access_token');
            if (!token) {
                setSubmitError('No hay sesión activa. Por favor, inicie sesión nuevamente.');
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
                return;
            }

            if (mode === 'create') {
                await usuarioService.createUsuarioAdmin(formData);
            } else {
                const updateData = {
                    nombre: formData.nombre.trim(),
                    apellido: formData.apellido.trim(),
                    email: formData.email.trim(),
                    is_admin: formData.is_admin,
                    password: formData.password
                };

                await usuarioService.updateUsuario(usuario.id_usuario, updateData);
            }

            onSave();
            onClose();
        } catch (error) {
            logger.error('UserFormModal handleSubmit error', error);
            setSubmitError(error.message || 'Error al conectar con el servidor');
        } finally {
            setIsSubmitting(false);
        }
    };

    const isEditMode = mode === 'edit';
    const modalTitle = isEditMode ? 'Editar Usuario' : 'Agregar Nuevo Usuario';
    const submitButtonText = isEditMode ? 'Guardar Cambios' : 'Crear Usuario';

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h2>{modalTitle}</h2>
                {submitError && <div className="error-message">{submitError}</div>}

                <form onSubmit={handleSubmit}>
                    <div>
                        <div className="form-group">
                            <label htmlFor="nombre">Nombre:</label>
                            <input
                                type="text"
                                id="nombre"
                                name="nombre"
                                value={formData.nombre}
                                onChange={handleChange}
                                placeholder="Nombre del usuario"
                                disabled={isSubmitting}
                                required
                            />
                            {validationErrors.nombre && <span className="error-text">{validationErrors.nombre}</span>}
                        </div>

                        <div className="form-group">
                            <label htmlFor="apellido">Apellido:</label>
                            <input
                                type="text"
                                id="apellido"
                                name="apellido"
                                value={formData.apellido}
                                onChange={handleChange}
                                placeholder="Apellido del usuario"
                                disabled={isSubmitting}
                                required
                            />
                            {validationErrors.apellido && <span className="error-text">{validationErrors.apellido}</span>}
                        </div>

                        <div className="form-group">
                            <label htmlFor="email">Email:</label>
                            <input
                                type="email"
                                id="email"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
                                placeholder="correo@ejemplo.com"
                                disabled={isSubmitting}
                                required
                            />
                            {validationErrors.email && <span className="error-text">{validationErrors.email}</span>}
                        </div>

                        <div className="form-group">
                            <label htmlFor="username">Nombre de usuario:</label>
                            <input
                                type="text"
                                id="username"
                                name="username"
                                value={formData.username}
                                onChange={handleChange}
                                placeholder="nombre_usuario"
                                disabled={isSubmitting || isEditMode}
                                required
                            />
                            {validationErrors.username && <span className="error-text">{validationErrors.username}</span>}
                        </div>

                        <div className="form-group">
                            <label htmlFor="password">Contraseña:</label>
                            <input
                                type="text"
                                id="password"
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                                placeholder={isEditMode ? "Dejar vacío para mantener la contraseña actual" : "Contraseña (mínimo 6 caracteres)"}
                                disabled={isSubmitting}
                                required={!isEditMode}
                            />
                            {validationErrors.password && <span className="error-text">{validationErrors.password}</span>}
                        </div>
                        
                        <div className="form-group">
                            <label htmlFor="confirm_password">Confirmar contraseña:</label>
                            <input
                                type="text"
                                id="confirm_password"
                                name="confirm_password"
                                value={formData.confirm_password}
                                onChange={handleChange}
                                placeholder={isEditMode ? "Dejar vacío para mantener la contraseña actual" : "Contraseña (mínimo 6 caracteres)"}
                                disabled={isSubmitting}
                                required={!isEditMode}
                            />
                            {validationErrors.confirm_password && <span className="error-text">{validationErrors.confirm_password}</span>}
                        </div>
                    </div>

                    <div className="form-buttons">
                        <button
                            type="submit"
                            className="btn-guardar"
                            disabled={isSubmitting}
                        >
                            {isSubmitting ? 'Guardando...' : submitButtonText}
                        </button>
                        <button
                            type="button"
                            className="btn-cancelar"
                            onClick={onClose}
                            disabled={isSubmitting}
                        >
                            Cancelar
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default UserFormModal;
