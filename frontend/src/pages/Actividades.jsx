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
        inscripciones,
        fetchById,
        fetchInscripciones,
        enrollInActividad,
        unenrollFromActividad,
        estaInscripto,
        deleteActividad
    } = useActividades(idUsuario);

    const [actividadesFiltradas, setActividadesFiltradas] = useState([]);
    const [totalItems, setTotalItems] = useState(0);
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
        dia: ""
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
        const { name, value } = e.target;
        setFiltros(prev => ({
            ...prev,
            [name]: value
        }));
        setPaginaActual(1);
    };

    const handleLimpiarFiltros = () => {
        setFiltros({
            busqueda: "",
            descripcion: "",
            dia: ""
        });
        setPaginaActual(1);
        // Search again with empty filters
        performSearch({
            id: "",
            busqueda: "",
            descripcion: "",
            page: 1
        });
    };

    const handleSearch = () => {
        performSearch(filtros, paginaActual);
    };

    const performSearch = async (currentFiltros, page = 1) => {
        setIsSearching(true);
        setSearchError(null);

        try {
            // Build search filters with pagination
            const searchFilters = {
                page: page,
                count: ITEMS_POR_PAGINA
            };
            if (currentFiltros.busqueda) searchFilters.titulo = currentFiltros.busqueda;
            if (currentFiltros.descripcion) searchFilters.descripcion = currentFiltros.descripcion;
            if (currentFiltros.dia) searchFilters.dia = currentFiltros.dia;

            logger.info('Searching with filters:', searchFilters);

            // Call Search API with pagination
            const response = await searchService.searchActivities(searchFilters);
            const results = response.results || [];

            logger.info('Search results:', results);

            setActividadesFiltradas(results);
            setTotalItems(response.total || 0);
            setPaginaActual(page);
        } catch (error) {
            logger.error('Error searching activities:', error);
            setSearchError('Error al buscar actividades. Por favor, intenta nuevamente.');
            setActividadesFiltradas([]);
            setTotalItems(0);
        } finally {
            setIsSearching(false);
        }
    };

    // Paginación (server-side)
    const totalPaginas = Math.ceil(totalItems / ITEMS_POR_PAGINA);

    const handlePrevPage = () => {
        const newPage = Math.max(paginaActual - 1, 1);
        if (newPage !== paginaActual) {
            performSearch(filtros, newPage);
        }
    };

    const handleNextPage = () => {
        const newPage = Math.min(paginaActual + 1, totalPaginas);
        if (newPage !== paginaActual) {
            performSearch(filtros, newPage);
        }
    };

    const handleEnroling = (actividad) => {
        if (!isLoggedIn) {
            setAlertDialog({
                title: 'No estás loggeado',
                message: 'Debes iniciar sesión para inscribirte',
                type: 'error'
            });
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
            // Refresh search results preserving current page
            performSearch(filtros, paginaActual);
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
            // Refresh search results preserving current page
            performSearch(filtros, paginaActual);
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

    const handleToggleExpand = async (actividad) => {
        if (!actividad || !actividad.id_actividad) {
            logger.error("Error: actividad sin ID", actividad);
            return;
        }

        if (!isLoggedIn) {
            setAlertDialog({
                title: '¿Estás listo para empezar?',
                message: 'Inicia sesión para ver mas',
                type: 'info'
            });
            return
        }

        try {
            // Fetch complete activity data from API
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

    const handleEditar = async (actividad) => {
        setExpandedActividad(null); // Cerramos el detalle expandido

        if (!actividad || !actividad.id_actividad) {
            logger.error("Error: actividad sin ID para editar", actividad);
            return;
        }

        try {
            // Fetch complete activity data from API
            const actividadCompleta = await fetchById(actividad.id_actividad);
            setActividadEditar(actividadCompleta);
        } catch (error) {
            logger.error('Error fetching activity for edit:', error);
            setAlertDialog({
                title: 'Error',
                message: 'No se pudo cargar los detalles de la actividad para editar',
                type: 'error'
            });
        }
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
        // Refresh search results preserving current page
        performSearch(filtros, paginaActual);
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
                        {actividadesFiltradas.map((actividad) => (
                            <ActivityCard
                                key={actividad.id_actividad}
                                actividad={actividad}
                                isLoggedIn={isLoggedIn}
                                isAdmin={isAdmin}
                                estaInscripto={estaInscripto}
                                onToggleExpand={handleToggleExpand}
                                onEditar={handleEditar}
                                onEliminar={handleEliminar}
                                onEnroling={handleEnroling}
                                onUnenrolling={handleUnenrolling}
                            />
                        ))}
                    </div>

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
                    inscriptionsEdit={true}
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
