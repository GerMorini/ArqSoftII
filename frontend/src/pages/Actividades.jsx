import React, { useState, useEffect } from "react";
import ActivityFormModal from '../components/ActivityFormModal';
import ActivityCard from '../components/ActivityCard';
import ActivityCardExpanded from '../components/ActivityCardExpanded';
import SearchBar from '../components/SearchBar';
import ConfirmDialog from '../components/ConfirmDialog';
import AlertDialog from '../components/AlertDialog';
import "../styles/Actividades.css";
import { useNavigate } from "react-router-dom";
import { useActividades } from '../hooks/useActividades';
import useCurrentUser from '../hooks/useCurrentUser';
import logger from '../utils/logger';
import searchService from '../services/searchService';

const Actividades = () => {
    const navigate = useNavigate();
    const { isLoggedIn, isAdmin, userId: idUsuario } = useCurrentUser();

    const {
        actividades,
        inscripciones,
        loading,
        error,
        fetchActividades,
        fetchInscripciones,
        enrollInActividad,
        unenrollFromActividad,
        estaInscripto,
        deleteActividad
    } = useActividades(idUsuario);

    const [actividadesFiltradas, setActividadesFiltradas] = useState([]);
    const [actividadEditar, setActividadEditar] = useState(null);
    const [expandedActividad, setExpandedActividad] = useState(null);
    const [actividadAEliminar, setActividadAEliminar] = useState(null);
    const [actividadADesincribir, setActividadADesincribir] = useState(null);
    const [actividadAInscribir, setActividadAInscribir] = useState(null);
    const [alertDialog, setAlertDialog] = useState(null);
    const [isSearching, setIsSearching] = useState(false);
    const [searchError, setSearchError] = useState(null);
    const [filtros, setFiltros] = useState({
        busqueda: "",
        descripcion: "",
        dia: "",
        soloInscripto: false
    });
    const [paginaActual, setPaginaActual] = useState(1);
    const ITEMS_POR_PAGINA = 9;

    useEffect(() => {
        // Initial load with empty search (shows all)
        handleSearch();
        // Solo cargar inscripciones si el usuario está loggeado y NO es admin
        if (isLoggedIn && !isAdmin && idUsuario) {
            fetchInscripciones(idUsuario);
        }
    }, []);

    const handleFiltroChange = (e) => {
        const { name, value, checked, type } = e.target;
        setFiltros(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));
        setPaginaActual(1);
    };

    const handleLimpiarFiltros = () => {
        setFiltros({
            busqueda: "",
            descripcion: "",
            dia: "",
            soloInscripto: false
        });
        setPaginaActual(1);
        // Search again with empty filters
        performSearch({
            busqueda: "",
            descripcion: "",
            dia: "",
            soloInscripto: false
        });
    };

    const handleSearch = () => {
        performSearch(filtros);
    };

    const performSearch = async (currentFiltros) => {
        setIsSearching(true);
        setSearchError(null);

        try {
            // Build search filters
            const searchFilters = {};
            if (currentFiltros.busqueda) searchFilters.titulo = currentFiltros.busqueda;
            if (currentFiltros.descripcion) searchFilters.descripcion = currentFiltros.descripcion;
            if (currentFiltros.dia) searchFilters.dia = currentFiltros.dia;

            logger.info('Searching with filters:', searchFilters);

            // Call Search API
            const response = await searchService.searchActivities(searchFilters);
            let results = response.results || [];

            logger.info('Search results:', results);

            // Apply client-side filter for "soloInscripto"
            if (currentFiltros.soloInscripto && inscripciones.length > 0) {
                results = results.filter(actividad =>
                    inscripciones.includes(actividad.id_actividad)
                );
            }

            setActividadesFiltradas(results);
            setPaginaActual(1);
        } catch (error) {
            logger.error('Error searching activities:', error);
            setSearchError('Error al buscar actividades. Por favor, intenta nuevamente.');
            setActividadesFiltradas([]);
        } finally {
            setIsSearching(false);
        }
    };

    // Paginación
    const totalPaginas = Math.ceil(actividadesFiltradas.length / ITEMS_POR_PAGINA);
    const inicio = (paginaActual - 1) * ITEMS_POR_PAGINA;
    const actividadesPaginadas = actividadesFiltradas.slice(inicio, inicio + ITEMS_POR_PAGINA);

    const handlePrevPage = () => {
        setPaginaActual(prev => Math.max(prev - 1, 1));
    };

    const handleNextPage = () => {
        setPaginaActual(prev => Math.min(prev + 1, totalPaginas));
    };

    const handleEnroling = (actividad) => {
        if (!isLoggedIn) {
            navigate("/login");
            return;
        }
        setActividadAInscribir(actividad);
    };

    const handleConfirmEnroll = async () => {
        try {
            await enrollInActividad(idUsuario, actividadAInscribir.id_actividad);
            logger.info('Inscripción completada, inscripciones actuales:', inscripciones);
            setAlertDialog({
                title: 'Inscripción exitosa',
                message: '¡Te has inscripto a la actividad!',
                type: 'success'
            });
            // Refresh search results
            handleSearch();
        } catch (error) {
            logger.error('handleConfirmEnroll error', error);
            setAlertDialog({
                title: 'Error al inscribirse',
                message: error.message || "Error al inscribirse en la actividad",
                type: 'error'
            });
        }
        setActividadAInscribir(null);
    };

    const handleCancelEnroll = () => {
        setActividadAInscribir(null);
    };

    const handleConfirmUnenroll = async () => {
        try {
            await unenrollFromActividad(idUsuario, actividadADesincribir.id_actividad);
            logger.info('Desinscripción completada, inscripciones actuales:', inscripciones);
            setAlertDialog({
                title: 'Desinscripción exitosa',
                message: 'Te has desincripto de la actividad',
                type: 'success'
            });
            // Refresh search results
            handleSearch();
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

    const handleUnenrolling = (actividad) => {
        setActividadADesincribir(actividad);
    };

    const handleEditar = (actividad) => {
        setExpandedActividad(null); // Cerramos el detalle expandido
        setActividadEditar(actividad);
    };

    const handleCloseModal = () => {
        setActividadEditar(null);
    };

    const handleSaveEdit = () => {
        setAlertDialog({
            title: 'Actividad actualizada',
            message: 'La actividad se ha actualizado exitosamente',
            type: 'success'
        });
        // Refresh search results
        handleSearch();
    };

    const handleEliminar = (actividad) => {
        if (!actividad.id_actividad) {
            logger.error("Error: La actividad no tiene ID", actividad);
            alert('Error: No se puede eliminar la actividad porque no tiene ID');
            return;
        }
        setActividadAEliminar(actividad);
    };

    const handleConfirmDelete = async () => {
        try {
            await deleteActividad(actividadAEliminar.id_actividad);
            setAlertDialog({
                title: 'Actividad eliminada',
                message: 'La actividad se ha eliminado exitosamente',
                type: 'success'
            });
            setActividadAEliminar(null);
            // Refresh search results
            handleSearch();
        } catch (error) {
            logger.error("handleConfirmDelete error", error);
            setAlertDialog({
                title: 'Error al eliminar',
                message: 'No se pudo eliminar la actividad',
                type: 'error'
            });
            setActividadAEliminar(null);
        }
    };

    const handleCancelDelete = () => {
        setActividadAEliminar(null);
    };

    return (
        <div className="actividades-container">
            <SearchBar
                filtros={filtros}
                onFiltroChange={handleFiltroChange}
                onLimpiar={handleLimpiarFiltros}
                onSearch={handleSearch}
                mostrarToggle={isLoggedIn && !isAdmin}
                soloInscriptoDisabled={false}
                isSearching={isSearching}
            />

            {searchError && (
                <div className="mensaje-error" style={{
                    padding: '1rem',
                    backgroundColor: '#ffebee',
                    color: '#c62828',
                    borderRadius: '8px',
                    marginBottom: '1rem',
                    textAlign: 'center'
                }}>
                    {searchError}
                </div>
            )}

            {isSearching ? (
                <div className="mensaje-no-actividades">
                    Buscando actividades...
                </div>
            ) : actividadesFiltradas.length === 0 ? (
                <div className="mensaje-no-actividades">
                    No se encontraron actividades.
                </div>
            ) : (
                <>
                    <div className="actividades-grid">
                        {actividadesPaginadas.map((actividad) => (
                            <ActivityCard
                                key={actividad.id_actividad}
                                actividad={actividad}
                                isLoggedIn={isLoggedIn}
                                isAdmin={isAdmin}
                                estaInscripto={estaInscripto}
                                onToggleExpand={setExpandedActividad}
                                onEditar={handleEditar}
                                onEliminar={handleEliminar}
                                onEnroling={handleEnroling}
                                onUnenrolling={handleUnenrolling}
                            />
                        ))}
                    </div>

                    {actividadesFiltradas.length > ITEMS_POR_PAGINA && (
                        <div className="pagination-controls">
                            <span style={{ padding: '0.5rem 1rem', color: '#2c3e50', fontWeight: '500' }}>
                                Página {paginaActual} de {totalPaginas}
                            </span>
                            <button
                                className="pagination-button"
                                onClick={handlePrevPage}
                                disabled={paginaActual === 1}
                                aria-label="Página anterior"
                            >
                                ← Anterior
                            </button>
                            <button
                                className="pagination-button"
                                onClick={handleNextPage}
                                disabled={paginaActual === totalPaginas}
                                aria-label="Página siguiente"
                            >
                                Siguiente →
                            </button>
                        </div>
                    )}
                </>
            )}

            {expandedActividad && (
                <ActivityCardExpanded
                    actividad={expandedActividad}
                    onClose={() => setExpandedActividad(null)}
                />
            )}

            {actividadEditar && (
                <ActivityFormModal
                    mode="edit"
                    actividad={actividadEditar}
                    onClose={handleCloseModal}
                    onSave={handleSaveEdit}
                />
            )}

            {actividadAInscribir && (
                <ConfirmDialog
                    title="Confirmar Inscripción"
                    message="¿Deseas inscribirse a esta actividad?"
                    details={`Se inscribirá en: "${actividadAInscribir.titulo}"`}
                    confirmText="Inscribirse"
                    cancelText="Cancelar"
                    isDangerous={false}
                    onConfirm={handleConfirmEnroll}
                    onCancel={handleCancelEnroll}
                />
            )}

            {actividadAEliminar && (
                <ConfirmDialog
                    title="Eliminar Actividad"
                    message="¿Estás seguro de que deseas eliminar esta actividad? Esta acción no se puede deshacer."
                    details={`Se eliminará: "${actividadAEliminar.titulo}"`}
                    confirmText="Eliminar"
                    cancelText="Cancelar"
                    isDangerous={true}
                    onConfirm={handleConfirmDelete}
                    onCancel={handleCancelDelete}
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

export default Actividades;