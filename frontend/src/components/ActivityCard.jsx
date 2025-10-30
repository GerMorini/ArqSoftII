import React from 'react';

const ActivityCard = ({
    actividad,
    isLoggedIn,
    isAdmin,
    estaInscripto,
    onToggleExpand,
    onEditar,
    onEliminar,
    onEnroling,
    onUnenrolling
}) => {
    return (
        <div className="actividad-card">
            <h3>{actividad.titulo}</h3>
            <div className="actividad-info-basic">
                <p>Instructor: {actividad.instructor || "No especificado"}</p>
                <p>
                    Horario: {actividad.hora_inicio} a {actividad.hora_fin}
                </p>
            </div>

            <div className="card-actions">
                {isLoggedIn && (
                    <>
                        {isAdmin ? (
                            <>
                                <button
                                    className="card-edit-button"
                                    onClick={() => onEditar(actividad)}
                                    title="Editar"
                                >
                                    ✏️
                                </button>
                                <button
                                    className="card-delete-button"
                                    onClick={() => onEliminar(actividad)}
                                    title="Eliminar"
                                >
                                    🗑️
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
                    onClick={() => onToggleExpand(actividad)}
                >
                    Ver más 🔽
                </button>
            </div>
        </div>
    );
};

export default ActivityCard;