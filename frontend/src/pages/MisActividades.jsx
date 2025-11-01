import React, { useState, useEffect } from 'react';
import ActivityCard from '../components/ActivityCard';
import ActivityCardExpanded from '../components/ActivityCardExpanded';
import ConfirmDialog from '../components/ConfirmDialog';
import AlertDialog from '../components/AlertDialog';
import '../styles/MisActividades.css';
import gymPortada from '../assets/login/gimnasio1.jpeg';
import useCurrentUser from '../hooks/useCurrentUser';
import { useActividades } from '../hooks/useActividades';
import { actividadService } from '../services/actividadService';
import logger from '../utils/logger';

const MisActividades = () => {
    const { isLoggedIn, userId: idUsuario } = useCurrentUser();
    const {
        fetchById,
        unenrollFromActividad,
        estaInscripto
    } = useActividades(idUsuario);

    const [misActividades, setMisActividades] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [expandedActividad, setExpandedActividad] = useState(null);
    const [actividadADesincribir, setActividadADesincribir] = useState(null);
    const [alertDialog, setAlertDialog] = useState(null);

    useEffect(() => {
        if (isLoggedIn && idUsuario) {
            loadMisActividades();
        }
    }, [isLoggedIn, idUsuario]);

    const loadMisActividades = async () => {
        setIsLoading(true);
        setError(null);
        try {
            logger.info('Cargando actividades inscritas para usuario:', idUsuario);
            const activities = await actividadService.getInscribedActivitiesData(idUsuario);
            setMisActividades(activities);
            logger.info('Actividades inscritas cargadas:', activities.length);
        } catch (err) {
            logger.error('Error al cargar actividades inscritas:', err);
            setError('Error al cargar tus actividades. Por favor, intenta nuevamente.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleToggleExpand = async (actividad) => {
        if (!actividad || !actividad.id_actividad) {
            logger.error("Error: actividad sin ID", actividad);
            return;
        }

        try {
            const actividadCompleta = await fetchById(actividad.id_actividad);
            setExpandedActividad(actividadCompleta);
        } catch (error) {
            logger.error('Error fetching activity details:', error);
            setAlertDialog({
                title: 'Error',
                message: 'No se pudo cargar los detalles de la actividad',
                type: 'error'
            });
        }
    };

    const handleUnenrolling = (actividad) => {
        setActividadADesincribir(actividad);
    };

    const handleConfirmUnenroll = async () => {
        try {
            await unenrollFromActividad(idUsuario, actividadADesincribir.id_actividad);
            logger.info('Desinscripción completada');
            setAlertDialog({
                title: 'Desinscripción exitosa',
                message: 'Te has desincripto de la actividad',
                type: 'success'
            });
            // Reload activities
            loadMisActividades();
        } catch (error) {
            logger.error("handleConfirmUnenroll error", error);
            setAlertDialog({
                title: 'Error al desincribirse',
                message: 'No se pudo desincribir de la actividad',
                type: 'error'
            });
        }
        setActividadADesincribir(null);
    };

    const handleCancelUnenroll = () => {
        setActividadADesincribir(null);
    };

    return (
        <div className="mis-actividades-page">
            <img
                className="background-gym"
                src={gymPortada}
                alt="Gimnasio background"
            />

            <div className="mis-actividades-content">
                <div className="page-header">
                    <h1 className="page-title">Mis Actividades</h1>
                    <p className="page-subtitle">
                        Aquí puedes ver todas las actividades en las que estás inscrito
                    </p>
                </div>

                {error && (
                    <div className="error-message">
                        {error}
                    </div>
                )}

                {isLoading ? (
                    <div className="loading-message">
                        Cargando tus actividades...
                    </div>
                ) : misActividades.length === 0 ? (
                    <div className="empty-message">
                        <div className="empty-icon">📋</div>
                        <h2>No tienes actividades inscritas</h2>
                        <p>Explora nuestras actividades disponibles y comienza tu entrenamiento</p>
                    </div>
                ) : (
                    <div className="actividades-grid">
                        {misActividades.map((actividad) => (
                            <ActivityCard
                                key={actividad.id_actividad}
                                actividad={actividad}
                                isLoggedIn={isLoggedIn}
                                isAdmin={false}
                                estaInscripto={(param) => true}
                                onToggleExpand={handleToggleExpand}
                                onUnenrolling={handleUnenrolling}
                            />
                        ))}
                    </div>
                )}
            </div>

            {expandedActividad && (
                <ActivityCardExpanded
                    actividad={expandedActividad}
                    onClose={() => setExpandedActividad(null)}
                />
            )}

            {actividadADesincribir && (
                <ConfirmDialog
                    title="Desincribirse"
                    message="¿Estás seguro de que deseas desincribirse de esta actividad?"
                    details={`Se desincribirá de: "${actividadADesincribir.titulo}"`}
                    confirmText="Desincribirse"
                    cancelText="Cancelar"
                    isDangerous={false}
                    onConfirm={handleConfirmUnenroll}
                    onCancel={handleCancelUnenroll}
                />
            )}

            {alertDialog && (
                <AlertDialog
                    title={alertDialog.title}
                    message={alertDialog.message}
                    type={alertDialog.type}
                    onClose={() => setAlertDialog(null)}
                />
            )}
        </div>
    );
};

export default MisActividades;
