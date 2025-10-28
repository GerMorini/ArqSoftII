import { useState, useMemo, useEffect } from 'react';
import ActivityFormModal from './ActivityFormModal';
import ConfirmDialog from './ConfirmDialog';
import AlertDialog from './AlertDialog';
import FilterBar from './FilterBar';
import '../styles/AdminPanel.css';
import '../styles/FilterBar.css';
import { useActividades } from '../hooks/useActividades';
import logger from '../utils/logger';

const AdminActivitiesTab = () => {
    const [actividadEditar, setActividadEditar] = useState(null);
    const [mostrarAgregarModal, setMostrarAgregarModal] = useState(false);
    const [actividadAEliminar, setActividadAEliminar] = useState(null);
    const [filtros, setFiltros] = useState({
        busqueda: '',
        descripcion: '',
        dia: ''
    });
    const [ordenamiento, setOrdenamiento] = useState({
        campo: null,
        direccion: 'asc'
    });
    const [paginaActual, setPaginaActual] = useState(1);
    const ITEMS_POR_PAGINA = 5;
    const [alertDialog, setAlertDialog] = useState(null);

    const {
        actividades,
        fetchActividades,
        deleteActividad,
    } = useActividades();

    useEffect(() => {
        fetchActividades();
    }, [fetchActividades]);

    const handleEditar = (actividad) => {
        setActividadEditar(actividad);
    };

    const handleCloseModal = () => {
        setActividadEditar(null);
        setMostrarAgregarModal(false);
    };

    const handleSaveEdit = () => {
        const isCreating = mostrarAgregarModal;
        const title = isCreating ? 'Actividad creada' : 'Actividad actualizada';
        const message = isCreating
            ? 'La actividad se ha creado exitosamente'
            : 'La actividad se ha actualizado exitosamente';
        setAlertDialog({ title, message, type: 'success' });
        fetchActividades();
        handleCloseModal();
        setOrdenamiento({ campo: null, direccion: 'asc' });
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
            fetchActividades();
            setAlertDialog({
                title: 'Actividad eliminada',
                message: 'La actividad se ha eliminado exitosamente',
                type: 'success'
            });
            setActividadAEliminar(null);
            setOrdenamiento({ campo: null, direccion: 'asc' });
            // setPaginaActual(1);
        } catch (error) {
            logger.error("handleConfirmDelete error", error);
            setAlertDialog({
                title: 'Error al eliminar',
                message: 'No se pudo eliminar la actividad. Por favor, intenta de nuevo más tarde.',
                type: 'error'
            });
            setActividadAEliminar(null);
        }
    };

    const handleCancelDelete = () => {
        setActividadAEliminar(null);
    };

    const handleFiltroChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFiltros(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));
        setPaginaActual(1);
    };

    const handleLimpiarFiltros = () => {
        setFiltros({
            busqueda: '',
            descripcion: '',
            dia: ''
        });
        setPaginaActual(1);
    };

    const handleOrdenarPor = (campo) => {
        setOrdenamiento(prev => ({
            campo,
            direccion: prev.campo === campo && prev.direccion === 'asc' ? 'desc' : 'asc'
        }));
    };

    const actividadesFiltradas = useMemo(() => {
        let resultado = actividades.filter(actividad => {
            const coincideBusqueda = actividad.titulo.toLowerCase().includes(filtros.busqueda.toLowerCase());
            const coincideDescripcion = actividad.descripcion.toLowerCase().includes(filtros.descripcion.toLowerCase());
            const coincideDia = filtros.dia === '' || actividad.dia === filtros.dia;

            return coincideBusqueda && coincideDescripcion && coincideDia;
        });

        if (ordenamiento.campo) {
            resultado.sort((a, b) => {
                let valorA = a[ordenamiento.campo];
                let valorB = b[ordenamiento.campo];

                if (typeof valorA === 'string') valorA = valorA.toLowerCase();
                if (typeof valorB === 'string') valorB = valorB.toLowerCase();

                if (valorA < valorB) return ordenamiento.direccion === 'asc' ? -1 : 1;
                if (valorA > valorB) return ordenamiento.direccion === 'asc' ? 1 : -1;
                return 0;
            });
        }

        return resultado;
    }, [actividades, filtros, ordenamiento]);

    const totalPaginas = Math.ceil(actividadesFiltradas.length / ITEMS_POR_PAGINA);
    const inicio = (paginaActual - 1) * ITEMS_POR_PAGINA;
    const actividadesPaginadas = actividadesFiltradas.slice(inicio, inicio + ITEMS_POR_PAGINA);

    return (
        <>
            <div className="admin-tabs-header">
                <button
                    className="btn-agregar"
                    onClick={() => setMostrarAgregarModal(true)}
                >
                    <span>+</span>
                    Agregar Actividad
                </button>
            </div>

            <div className="admin-filters">
                <FilterBar
                    filtros={filtros}
                    onFiltroChange={handleFiltroChange}
                    onLimpiar={handleLimpiarFiltros}
                    mostrarToggle={false}
                />
            </div>

            <div className="admin-table-container">
                <table className="admin-table">
                    <thead>
                        <tr>
                            <th
                                onClick={() => handleOrdenarPor('titulo')}
                                className={`sortable ${ordenamiento.campo === 'titulo' ? `sorted-${ordenamiento.direccion}` : ''}`}
                            >
                                Título
                            </th>
                            <th
                                onClick={() => handleOrdenarPor('descripcion')}
                                className={`sortable ${ordenamiento.campo === 'descripcion' ? `sorted-${ordenamiento.direccion}` : ''}`}
                            >
                                Descripción
                            </th>
                            <th
                                onClick={() => handleOrdenarPor('instructor')}
                                className={`sortable ${ordenamiento.campo === 'instructor' ? `sorted-${ordenamiento.direccion}` : ''}`}
                            >
                                Instructor
                            </th>
                            <th
                                onClick={() => handleOrdenarPor('dia')}
                                className={`sortable ${ordenamiento.campo === 'dia' ? `sorted-${ordenamiento.direccion}` : ''}`}
                            >
                                Día
                            </th>
                            <th
                                onClick={() => handleOrdenarPor('hora_inicio')}
                                className={`sortable ${ordenamiento.campo === 'hora_inicio' ? `sorted-${ordenamiento.direccion}` : ''}`}
                            >
                                Horario
                            </th>
                            <th
                                onClick={() => handleOrdenarPor('cupo')}
                                className={`sortable ${ordenamiento.campo === 'cupo' ? `sorted-${ordenamiento.direccion}` : ''}`}
                            >
                                Inscriptos / Cupo
                            </th>
                            <th className="acciones-column">Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {actividadesPaginadas.map((actividad) => {
                            const lugaresDisponibles = actividad.lugares_disponibles;
                            const estaLleno = lugaresDisponibles <= 0;

                            return (
                                <tr key={actividad.id_actividad}>
                                    <td>{actividad.titulo}</td>
                                    <td>{actividad.descripcion}</td>
                                    <td>{actividad.instructor}</td>
                                    <td>{actividad.dia}</td>
                                    <td>{actividad.hora_inicio} - {actividad.hora_fin}</td>
                                    <td className="cupo-cell">
                                        <span className={estaLleno ? 'cupo-lleno' : 'cupo-disponible'}>
                                            {actividad.cupo - lugaresDisponibles} / {actividad.cupo}
                                        </span>
                                    </td>
                                    <td className="acciones-column">
                                        <button
                                            className="action-button edit-button"
                                            onClick={() => handleEditar(actividad)}
                                            title="Editar"
                                        >
                                            ✏️
                                        </button>
                                        <button
                                            className="action-button delete-button"
                                            onClick={() => handleEliminar(actividad)}
                                            title="Eliminar"
                                        >
                                            🗑️
                                        </button>
                                    </td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>

                {totalPaginas > 1 && (
                    <div className="pagination-container">
                        <span className="pagination-info">
                            Mostrando {inicio + 1} a {Math.min(inicio + ITEMS_POR_PAGINA, actividadesFiltradas.length)} de {actividadesFiltradas.length} actividades
                        </span>
                        <div className="pagination-controls">
                            <button
                                className="pagination-btn"
                                onClick={() => setPaginaActual(prev => Math.max(prev - 1, 1))}
                                disabled={paginaActual === 1}
                            >
                                ← Anterior
                            </button>
                            <span style={{ padding: '0.5rem 1rem', color: '#2c3e50', fontWeight: '500' }}>
                                Página {paginaActual} de {totalPaginas}
                            </span>
                            <button
                                className="pagination-btn"
                                onClick={() => setPaginaActual(prev => Math.min(prev + 1, totalPaginas))}
                                disabled={paginaActual === totalPaginas}
                            >
                                Siguiente →
                            </button>
                        </div>
                    </div>
                )}
            </div>

            {actividadEditar && (
                <ActivityFormModal
                    mode="edit"
                    actividad={actividadEditar}
                    onClose={handleCloseModal}
                    onSave={handleSaveEdit}
                />
            )}

            {mostrarAgregarModal && (
                <ActivityFormModal
                    mode="create"
                    onClose={handleCloseModal}
                    onSave={handleSaveEdit}
                />
            )}

            {actividadAEliminar && (
                <ConfirmDialog
                    title="Eliminar Actividad"
                    message="¿Estás seguro de que deseas eliminar esta actividad? Se eliminarán también todas las inscripciones asociadas. Esta acción no se puede deshacer."
                    details={`Se eliminará: "${actividadAEliminar.titulo}"`}
                    confirmText="Eliminar"
                    cancelText="Cancelar"
                    isDangerous={true}
                    onConfirm={handleConfirmDelete}
                    onCancel={handleCancelDelete}
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
        </>
    );
};

export default AdminActivitiesTab;
