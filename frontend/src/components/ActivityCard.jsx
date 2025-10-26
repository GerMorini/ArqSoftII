import React from 'react';
import { useEscapeKey } from '../hooks/useEscapeKey';
import logger from '../utils/logger';

const ActivityCard = ({
    actividad,
    isExpanded,
    isLoggedIn,
    isAdmin,
    estaInscripto,
    onToggleExpand,
    onEditar,
    onEliminar,
    onEnroling,
    onUnenrolling
}) => {
    useEscapeKey(() => {
        if (isExpanded) {
            logger.debug('ActivityCard - Escape key pressed, closing');
            onToggleExpand(null);
        }
    });

    return (
        <div
            className={`actividad-card ${isExpanded ? 'expanded' : ''}`}
        >
            <h3>{actividad.titulo}</h3>
            <div className="actividad-info-basic">
                <p>Instructor: {actividad.instructor || "No especificado"}</p>
                <p>
                    Horario: {actividad.hora_inicio} a {actividad.hora_fin}
                </p>
            </div>

            {isExpanded && (
                <div className="actividad-info-expanded">
                    <div className="actividad-imagen">
                        <img
                            src={actividad.foto_url || "https://via.placeholder.com/300x200"}
                            alt={actividad.titulo}
                        />
                    </div>
                    <div className="actividad-detalles">
                        <p>{actividad.descripcion}</p>
                        <p>Día: {actividad.dia || "No especificado"}</p>
                        <p><b>Horario:</b> {actividad.hora_inicio} a {actividad.hora_fin}</p>
                        <p>Lugares disponibles: {actividad.lugares_disponibles}</p>
                    </div>
                </div>
            )}

            <div className="card-actions">
                {isLoggedIn && (
                    <>
                        {isAdmin ? (
                            <>
                                <button
                                    className="edit-button"
                                    onClick={() => onEditar(actividad)}
                                    title="Editar"
                                >
                                    <span>✏️</span>
                                    Editar
                                </button>
                                <button
                                    className="delete-button"
                                    onClick={() => onEliminar(actividad)}
                                    title="Eliminar"
                                >
                                    <span>🗑️</span>
                                    Eliminar
                                </button>
                            </>
                        ) : (
                            <button
                                className="inscripcion-button"
                                onClick={() =>
                                    estaInscripto(actividad.id_actividad) ?
                                        onUnenrolling(actividad) :
                                        onEnroling(actividad)
                                }
                            >
                                {estaInscripto(actividad.id_actividad) ? "Desinscribir ❌" : "Inscribir ✔️"}
                            </button>
                        )}
                    </>
                )}
                <button
                    className="ver-mas-button"
                    onClick={() => onToggleExpand(isExpanded ? null : actividad.id_actividad)}
                >
                    {isExpanded ? "Ver menos 🔼" : "Ver más 🔽"}
                </button>
            </div>
        </div>
    );
};

export default ActivityCard;