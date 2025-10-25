import React, { useState, useEffect } from "react";
import EditarActividadModal from '../components/EditarActividadModal';
import ActivityCard from '../components/ActivityCard';
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
        descripcion: "",
        dia: "",
        soloInscripto: false
    });
    const isLoggedIn = localStorage.getItem("isLoggedIn") === "true";
    const isAdmin = localStorage.getItem("isAdmin") === "true";
    const idUsuario = localStorage.getItem("idUsuario")
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
                setActividades(data.activities);
                setActividadesFiltradas(data.activities);
            }
        } catch (error) {
            console.error("Error al cargar actividades:", error);
        }
    };

    const fetchInscripciones = async () => {
        try {
            console.log(`ACTIVITIES_URL = ${ACTIVITIES_URL}`);
            const response = await fetch(`${ACTIVITIES_URL}/inscriptions/${idUsuario}`, {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('access_token')}`
                },
            });
            if (response.ok) {
                const resp = await response.json();
                const data = resp.inscripciones

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
        if (filtros.descripcion) {
            const descLower = filtros.descripcion.toLowerCase();
            actividadesFiltradas = actividadesFiltradas.filter(actividad =>
                actividad.descripcion.toLowerCase().includes(descLower)
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
            actividadesFiltradas = actividadesFiltradas.filter(actividad =>
                inscripciones.includes(actividad.id_actividad)
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
            const response = await fetch(`${ACTIVITIES_URL}/activities/${actividadId}/inscribir`, {
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("access_token")}`
                }
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
            const response = await fetch(`${ACTIVITIES_URL}/activities/${id_actividad}/desinscribir`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('access_token')}`
                }
            });

            if (response.status == 200) {
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
                const response = await fetch(`${ACTIVITIES_URL}/activities/${actividad.id_actividad}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('access_token')}`
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
        return inscripciones.includes(id_actividad)
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
                    name="descripcion"
                    placeholder="Descripción..."
                    value={filtros.descripcion}
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
                        <ActivityCard
                            key={actividad.id_actividad}
                            actividad={actividad}
                            isExpanded={expandedActividadId === actividad.id_actividad}
                            isLoggedIn={isLoggedIn}
                            isAdmin={isAdmin}
                            estaInscripto={estaInscripto}
                            onToggleExpand={setExpandedActividadId}
                            onEditar={handleEditar}
                            onEliminar={handleEliminar}
                            onEnroling={handleEnroling}
                            onUnenrolling={handleUnenrolling}
                        />
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