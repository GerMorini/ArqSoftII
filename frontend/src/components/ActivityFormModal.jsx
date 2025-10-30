import React, { useState, useEffect } from 'react';
import '../styles/ActivityFormModal.css';
import { useEscapeKey } from '../hooks/useEscapeKey';
import { DIAS_SEMANA } from '../constants/actividadConstants';
import { actividadService } from '../services/actividadService';
import logger from '../utils/logger';
import useCurrentUser from '../hooks/useCurrentUser';

const ActivityFormModal = ({ mode = 'create', actividad = null, onClose, onSave }) => {
    const [formData, setFormData] = useState({
        id_actividad: '',
        titulo: '',
        descripcion: '',
        cupo: '',
        dia: '',
        hora_inicio: '',
        hora_fin: '',
        foto_url: '',
        instructor: '',
        usuarios_inscritos_text: ''
    });
    const [submitError, setSubmitError] = useState('');
    const [validationErrors, setValidationErrors] = useState({});
    const [isSubmitting, setIsSubmitting] = useState(false);

    const currentUser = useCurrentUser();

    useEscapeKey(onClose);

    // Inicializar formulario según modo
    useEffect(() => {
        if (mode === 'edit' && actividad) {
            const actividadData = {
                id_actividad: actividad.id_actividad,
                titulo: actividad.titulo || '',
                descripcion: actividad.descripcion || '',
                cupo: actividad.cupo || '',
                dia: actividad.dia || '',
                hora_inicio: actividad.hora_inicio || '',
                hora_fin: actividad.hora_fin || '',
                foto_url: actividad.foto_url || '',
                instructor: actividad.instructor || ''
            };

            // If admin view includes usuarios_inscritos, expose as comma-separated text for editing
            if (currentUser.isAdmin && actividad.usuarios_inscritos) {
                actividadData.usuarios_inscritos_text = actividad.usuarios_inscritos.join(',');
            } else {
                actividadData.usuarios_inscritos_text = '';
            }

            setFormData(actividadData);
        } else {
            setFormData({
                id_actividad: '',
                titulo: '',
                descripcion: '',
                cupo: '',
                dia: '',
                hora_inicio: '',
                hora_fin: '',
                foto_url: '',
                instructor: '',
                usuarios_inscritos_text: ''
            });
        }
    }, [mode, actividad, currentUser]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
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
        const errors = actividadService.validateActividadForm(formData);
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

            const dataToSend = {
                ...formData,
                cupo: formData.cupo.toString()
            };

            // If admin edited the inscritos text, convert to array and include as usuarios_inscritos
            if (currentUser.isAdmin && formData.usuarios_inscritos_text !== undefined) {
                const raw = formData.usuarios_inscritos_text || '';
                const list = raw.split(',').map(s => s.trim()).filter(s => s !== '');
                dataToSend.usuarios_inscritos = list;
            }

            if (mode === 'create') {
                await actividadService.createActividad(dataToSend);
            } else {
                await actividadService.updateActividad(formData.id_actividad, dataToSend);
            }

            onSave();
            onClose();
        } catch (error) {
            logger.error('ActivityFormModal handleSubmit error', error);
            setSubmitError(error.message || 'Error al conectar con el servidor');
        } finally {
            setIsSubmitting(false);
        }
    };

    const isEditMode = mode === 'edit';
    const modalTitle = isEditMode ? 'Editar Actividad' : 'Agregar Nueva Actividad';
    const submitButtonText = isEditMode ? 'Guardar Cambios' : 'Crear Actividad';

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h2>{modalTitle}</h2>
                {submitError && <div className="error-message">{submitError}</div>}

                <form onSubmit={handleSubmit}>
                    <div>
                        {/* Columna Izquierda */}
                        <div className="form-group">
                            <label htmlFor="titulo">Título:</label>
                            <input
                                type="text"
                                id="titulo"
                                name="titulo"
                                value={formData.titulo}
                                onChange={handleChange}
                                placeholder="Nombre de la actividad"
                                disabled={isSubmitting}
                                required
                            />
                            {validationErrors.titulo && <span className="error-text">{validationErrors.titulo}</span>}
                        </div>

                        <div className="form-group">
                            <label htmlFor="instructor">Instructor:</label>
                            <input
                                type="text"
                                id="instructor"
                                name="instructor"
                                value={formData.instructor}
                                onChange={handleChange}
                                placeholder="Nombre del instructor"
                                disabled={isSubmitting}
                                required
                            />
                            {validationErrors.instructor && <span className="error-text">{validationErrors.instructor}</span>}
                        </div>

                        <div className="dia-cupo-container">
                            <div className="form-group">
                                <label htmlFor="dia">Día:</label>
                                <select
                                    id="dia"
                                    name="dia"
                                    value={formData.dia}
                                    onChange={handleChange}
                                    disabled={isSubmitting}
                                    required
                                >
                                    <option value="">Seleccione un día</option>
                                    {DIAS_SEMANA.map((dia) => (
                                        <option key={dia.value} value={dia.value}>
                                            {dia.label}
                                        </option>
                                    ))}
                                </select>
                                {validationErrors.dia && <span className="error-text">{validationErrors.dia}</span>}
                            </div>

                            <div className="form-group">
                                <label htmlFor="cupo">
                                    Cupo:
                                    {isEditMode && actividad && (
                                        <span className="inscriptos-info"> ({actividad.cupo - actividad.lugares_disponibles} inscriptos)</span>
                                    )}
                                </label>
                                <input
                                    type="number"
                                    id="cupo"
                                    name="cupo"
                                    value={formData.cupo}
                                    onChange={handleChange}
                                    placeholder="Cantidad de lugares"
                                    disabled={isSubmitting}
                                    required
                                    min="1"
                                />
                                {validationErrors.cupo && <span className="error-text">{validationErrors.cupo}</span>}
                            </div>
                        </div>

                        <div className="horarios-container">
                            <div className="form-group">
                                <label htmlFor="hora_inicio">Hora de inicio:</label>
                                <input
                                    type="time"
                                    id="hora_inicio"
                                    name="hora_inicio"
                                    value={formData.hora_inicio}
                                    onChange={handleChange}
                                    step="1800"
                                    disabled={isSubmitting}
                                    required
                                />
                                {validationErrors.hora_inicio && <span className="error-text">{validationErrors.hora_inicio}</span>}
                            </div>

                            <div className="form-group">
                                <label htmlFor="hora_fin">Hora de fin:</label>
                                <input
                                    type="time"
                                    id="hora_fin"
                                    name="hora_fin"
                                    value={formData.hora_fin}
                                    onChange={handleChange}
                                    step="1800"
                                    disabled={isSubmitting}
                                    required
                                />
                                {validationErrors.hora_fin && <span className="error-text">{validationErrors.hora_fin}</span>}
                            </div>
                        </div>

                        {/* Columna Derecha */}
                        <div className="form-group">
                            <label htmlFor="descripcion">Descripción:</label>
                            <textarea
                                id="descripcion"
                                name="descripcion"
                                value={formData.descripcion}
                                onChange={handleChange}
                                placeholder="Descripción de la actividad"
                                disabled={isSubmitting}
                                required
                            />
                            {validationErrors.descripcion && <span className="error-text">{validationErrors.descripcion}</span>}
                        </div>

                        <div className="form-group">
                            <label htmlFor="foto_url">URL de la foto:</label>
                            <input
                                type="text"
                                id="foto_url"
                                name="foto_url"
                                value={formData.foto_url}
                                onChange={handleChange}
                                placeholder="https://ejemplo.com/foto.jpg"
                                disabled={isSubmitting}
                            />
                            {validationErrors.foto_url && <span className="error-text">{validationErrors.foto_url}</span>}
                        </div>

                        {currentUser.isAdmin && (
                            <div className="form-group">
                                <label htmlFor="usuarios_inscritos_text">Usuarios inscritos (IDs, separados por comas):</label>
                                <textarea
                                    id="usuarios_inscritos_text"
                                    name="usuarios_inscritos_text"
                                    value={formData.usuarios_inscritos_text || ''}
                                    onChange={handleChange}
                                    placeholder="123, 456, 789"
                                    disabled={isSubmitting}
                                />
                                <small>Dejar vacío para no modificar. Para eliminar todos, enviar campo vacío.</small>
                            </div>
                        )}
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

export default ActivityFormModal;
