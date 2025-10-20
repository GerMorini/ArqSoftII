import React, { useState, useEffect } from "react";
import EditarActividadModal from '../components/EditarActividadModal';
import "../styles/Actividades.css";
import { useNavigate } from "react-router-dom";
import config from '../config/env';

const Actividades = () => {
    const [actividades, setActividades] = useState([]);
    const [actividadesFiltradas, setActividadesFiltradas] = useState([]);
    const [inscripciones, setInscripciones] = useState([]);
    const [actividadEditar, setActividadEditar] = useState(null);
    const [expandedActividadId, setExpandedActividadId] = useState(null);
    const [filtros, setFiltros] = useState({
        busqueda: "",
        categoria: "",
        dia: "",
        soloInscripto: false
    });
    const isLoggedIn = localStorage.getItem("isLoggedIn") === "true";
    const isAdmin = localStorage.getItem("isAdmin") === "true";
    const navigate = useNavigate();
    const ACTIVITIES_URL = config.ACTIVITIES_URL;

    useEffect(() => {
        fetchActividades();
        fetchInscripciones();
    }, []);

    useEffect(() => {
        filtrarActividades();
    }, [filtros, actividades]);

    const fetchActividades = async () => {
        try {
            console.log(`ACTIVITIES_URL = ${ACTIVITIES_URL}`);
            const response = await fetch(`${ACTIVITIES_URL}/activities`);
            if (response.ok) {
                const data = await response.json();
                console.log("Actividades cargadas:", data);
                setActividades(data);
                setActividadesFiltradas(data);
            }
        } catch (error) {
            console.error("Error al cargar actividades:", error);
        }
    };
    
    const fetchInscripciones = async () => {
        try {
            console.log(`ACTIVITIES_URL = ${ACTIVITIES_URL}`);
            const response = await fetch(`${ACTIVITIES_URL}/inscriptions`, {
                headers: {'Authorization': `Bearer ${localStorage.getItem('access_token')}`
            },
            });
            if (response.ok) {
                const resp = await response.json();
                const data = resp.filter(insc => insc.is_activa)
                
                console.log("Inscripciones cargadas:", data);
                setInscripciones(data);
            }
        } catch (error) {
            console.error("Error al cargar inscripciones:", error);
        }
    };

    const handleFiltroChange = (e) => {
        const { name, value } = e.target;
        setFiltros(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const filtrarActividades = () => {
        let actividadesFiltradas = [...actividades];

        // Filtrar por búsqueda (título o descripción)
        if (filtros.busqueda) {
            const busquedaLower = filtros.busqueda.toLowerCase();
            actividadesFiltradas = actividadesFiltradas.filter(actividad =>
                actividad.titulo.toLowerCase().includes(busquedaLower) ||
                actividad.descripcion.toLowerCase().includes(busquedaLower)
            );
        }

        // Filtrar por categoría (ahora como búsqueda de texto)
        if (filtros.categoria) {
            const categoriaLower = filtros.categoria.toLowerCase();
            actividadesFiltradas = actividadesFiltradas.filter(actividad =>
                actividad.categoria.toLowerCase().includes(categoriaLower)
            );
        }

        // Filtrar por día
        if (filtros.dia) {
            actividadesFiltradas = actividadesFiltradas.filter(actividad =>
                actividad.dia.toLowerCase() === filtros.dia.toLowerCase()
            );
        }

        // Filtrar solo inscripto
        if (filtros.soloInscripto) {
            const idsInscripto = inscripciones.filter(insc => insc.is_activa).map(insc => insc.id_actividad);
            actividadesFiltradas = actividadesFiltradas.filter(actividad =>
                idsInscripto.includes(actividad.id_actividad)
            );
        }

        setActividadesFiltradas(actividadesFiltradas);
    };

    const handleEnroling = async (actividadId) => {
        if (!isLoggedIn) {
            navigate("/login");
            return;
        }

        try {
            console.log(`ACTIVITIES_URL = ${ACTIVITIES_URL}`);
            const response = await fetch(`${ACTIVITIES_URL}/inscriptions`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${localStorage.getItem("access_token")}`,
                },
                body: JSON.stringify({
                    id: actividadId,
                }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || "Error al inscribirse en la actividad");
            }

            // Actualizar la lista de inscripciones
            fetchInscripciones();
            // Actualizar la lista de actividades para reflejar el cambio en los cupos
            fetchActividades();
            alert("¡Inscripción exitosa!");
        } catch (error) {
            console.error("Error al inscribirse:", error);
            alert(error.message);
        }
    };
    
    const handleUnenrolling = async (id_actividad) => {
        try {
            console.log(`ACTIVITIES_URL = ${ACTIVITIES_URL}`);
            const response = await fetch(`${ACTIVITIES_URL}/inscriptions`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
                },
                body: JSON.stringify({
                    id: parseInt(id_actividad)
                })
            });

            if (response.status == 204) {
                alert(`Desinscripto exitosamente`);
                fetchInscripciones();
            } else {
                alert(`Ups! algo salio mal, vuelve a intentarlo mas tarde`);
            }

            fetchActividades();
        } catch (error) {
            alert(`Ups! algo salio mal, vuelve a intentarlo mas tarde`);
            console.error("Error al desinscribir el usuario:", error);
        }
    };

    const handleEditar = (actividad) => {
        setExpandedActividadId(null); // Cerramos el detalle expandido
        setActividadEditar(actividad);
    };

    const handleCloseModal = () => {
        setActividadEditar(null);
    };

    const handleSaveEdit = () => {
        fetchActividades();
    };

    const handleEliminar = async (actividad) => {
        if (!actividad.id_actividad) {
            console.error("Error: La actividad no tiene ID", actividad);
            alert('Error: No se puede eliminar la actividad porque no tiene ID');
            return;
        }

        if (window.confirm('¿Estás seguro de que deseas eliminar esta actividad?')) {
            try {
                console.log("Intentando eliminar actividad con ID:", actividad.id_actividad);
                console.log(`ACTIVITIES_URL = ${ACTIVITIES_URL}`);
                const response = await fetch(`${ACTIVITIES_URL}/actividades/${actividad.id_actividad}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
                        'Content-Type': 'application/json'
                    }
                });

                if (response.ok) {
                    fetchActividades();
                    alert('Actividad eliminada con éxito');
                } else {
                    const errorData = await response.json().catch(() => ({}));
                    alert(errorData.message || 'Error al eliminar la actividad');
                }
            } catch (error) {
                console.error("Error al eliminar:", error);
                alert('Error al eliminar la actividad');
            }
        }
    };

    const estaInscripto = (id_actividad) => {
        return inscripciones.some(insc => 
            insc.id_actividad === id_actividad &&
            insc.is_activa
        )
    };

    const toggleExpand = (actividadId) => {
        setExpandedActividadId(expandedActividadId === actividadId ? null : actividadId);
    };

    return (
        <div className="actividades-container">
            {expandedActividadId && (
                <div className="actividades-modal-bg" onClick={() => setExpandedActividadId(null)} />
            )}
            <div className="filtros-container">
                <div className="search-wrapper">
                    <span className="search-icon">🔍</span>
                    <input
                        type="text"
                        name="busqueda"
                        placeholder="Buscar actividad..."
                        value={filtros.busqueda}
                        onChange={handleFiltroChange}
                        className="filtro-input"
                    />
                </div>
                <input
                    type="text"
                    name="categoria"
                    placeholder="Categoría..."
                    value={filtros.categoria}
                    onChange={handleFiltroChange}
                    className="filtro-input"
                />
                <select
                    name="dia"
                    value={filtros.dia}
                    onChange={handleFiltroChange}
                    className="filtro-select"
                >
                    <option value="">Día</option>
                    <option value="Lunes">Lunes</option>
                    <option value="Martes">Martes</option>
                    <option value="Miercoles">Miercoles</option>
                    <option value="Jueves">Jueves</option>
                    <option value="Viernes">Viernes</option>
                    <option value="Sabado">Sabado</option>
                    <option value="Domingo">Domingo</option>
                </select>
                {isLoggedIn && !isAdmin && (
                    <div className="toggle-wrapper">
                        <label className="toggle-label">
                            <input
                                type="checkbox"
                                name="soloInscripto"
                                checked={filtros.soloInscripto}
                                onChange={(e) => setFiltros(prev => ({
                                    ...prev,
                                    soloInscripto: e.target.checked
                                }))}
                                className="toggle-input"
                            />
                            <span className="toggle-slider"></span>
                            <span className="toggle-text">Solo inscriptas</span>
                        </label>
                    </div>
                )}
            </div>

            <div className="actividades-grid">
                {actividadesFiltradas.length === 0 ? (
                    <div className="mensaje-no-actividades">
                        No se encontraron actividades.
                    </div>
                ) : (
                    actividadesFiltradas.map((actividad) => (
                        <div 
                            className={`actividad-card ${expandedActividadId === actividad.id_actividad ? 'expanded' : ''}`} 
                            key={actividad.id_actividad}
                        >
                            <h3>{actividad.titulo}</h3>
                            <div className="actividad-info-basic">
                                <p>Instructor: {actividad.instructor || "No especificado"}</p>
                                <p>
                                    Horario: {actividad.hora_inicio} a {actividad.hora_fin}
                                </p>
                            </div>

                            {expandedActividadId === actividad.id_actividad && (
                                <div className="actividad-info-expanded">
                                    <div className="actividad-imagen">
                                        <img 
                                            src={actividad.foto_url || "https://via.placeholder.com/300x200"} 
                                            alt={actividad.titulo}
                                        />
                                    </div>
                                    <div className="actividad-detalles">
                                        <p>{actividad.descripcion}</p>
                                        <p>Categoría: {actividad.categoria || "No especificada"}</p>
                                        <p>Día: {actividad.dia || "No especificado"}</p>
                                        <p><b>Horario:</b> {actividad.hora_inicio} a {actividad.hora_fin}</p>
                                        <p>Cupo total: {actividad.cupo} | Lugares disponibles: {actividad.lugares}</p>
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
                                                    onClick={() => handleEditar(actividad)}
                                                    title="Editar"
                                                >
                                                    <span>✏️</span>
                                                    Editar
                                                </button>
                                                <button
                                                    className="delete-button"
                                                    onClick={() => handleEliminar(actividad)}
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
                                                        handleUnenrolling(actividad.id_actividad) :
                                                        handleEnroling(actividad.id_actividad)
                                                }
                                            >
                                                {estaInscripto(actividad.id_actividad) ? "Desinscribir ❌" : "Inscribir ✔️"}
                                            </button>
                                        )}
                                    </>
                                )}
                                <button
                                    className="ver-mas-button"
                                    onClick={() => toggleExpand(actividad.id_actividad)}
                                >
                                    {expandedActividadId === actividad.id_actividad ? "Ver menos 🔼" : "Ver más 🔽"}
                                </button>
                            </div>
                        </div>
                    ))
                )}
            </div>

            {actividadEditar && (
                <EditarActividadModal
                    actividad={actividadEditar}
                    onClose={handleCloseModal}
                    onSave={handleSaveEdit}
                />
            )}
        </div>
    );
};

export default Actividades;
