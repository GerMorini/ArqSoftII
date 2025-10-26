import { useState, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import ActivityFormModal from '../components/ActivityFormModal';
import UserFormModal from '../components/UserFormModal';
import ConfirmDialog from '../components/ConfirmDialog';
import AlertDialog from '../components/AlertDialog';
import FilterBar from '../components/FilterBar';
import '../styles/AdminPanel.css';
import '../styles/FilterBar.css';
import { useActividades } from '../hooks/useActividades';
import { useUsuarios } from '../hooks/useUsuarios';
import { usuarioService } from '../services/usuarioService';
import logger from '../utils/logger';

const AdminPanel = () => {
    const [tabActiva, setTabActiva] = useState('actividades');

    // Estados de actividades
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
    const ITEMS_POR_PAGINA = 10;

    // Estados de usuarios
    const [usuarioEditar, setUsuarioEditar] = useState(null);
    const [mostrarAgregarUsuarioModal, setMostrarAgregarUsuarioModal] = useState(false);
    const [usuarioAEliminar, setUsuarioAEliminar] = useState(null);
    const [filtrosUsuarios, setFiltrosUsuarios] = useState({
        busqueda: '',
        email: '',
        username: '',
        isAdmin: ''
    });
    const [ordenamientoUsuarios, setOrdenamientoUsuarios] = useState({
        campo: null,
        direccion: 'asc'
    });
    const [paginaActualUsuarios, setPaginaActualUsuarios] = useState(1);

    const [alertDialog, setAlertDialog] = useState(null);
    const navigate = useNavigate();

    const {
        actividades,
        fetchActividades,
        deleteActividad,
    } = useActividades();

    const {
        usuarios,
        fetchUsuarios,
        deleteUsuario,
        updateUsuario,
    } = useUsuarios();

    useEffect(() => {
        const isAdmin = localStorage.getItem("isAdmin") === "true";
        if (!isAdmin) {
            navigate('/');
            return;
        }
        fetchActividades();
        fetchUsuarios();
    }, [navigate]);

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
        // Reiniciar ordenamiento después de actualizar
        setOrdenamiento({ campo: null, direccion: 'asc' });
        setPaginaActual(1);
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
            // Reiniciar ordenamiento después de eliminar
            setOrdenamiento({ campo: null, direccion: 'asc' });
            setPaginaActual(1);
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

    // Manejar cambios en filtros
    const handleFiltroChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFiltros(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));
        setPaginaActual(1);
    };

    // Limpiar filtros
    const handleLimpiarFiltros = () => {
        setFiltros({
            busqueda: '',
            descripcion: '',
            dia: ''
        });
        setPaginaActual(1);
    };

    // Manejar ordenamiento
    const handleOrdenarPor = (campo) => {
        setOrdenamiento(prev => ({
            campo,
            direccion: prev.campo === campo && prev.direccion === 'asc' ? 'desc' : 'asc'
        }));
        setPaginaActual(1);
    };

    // Filtrar y ordenar actividades
    const actividadesFiltradas = useMemo(() => {
        let resultado = actividades.filter(actividad => {
            const coincideBusqueda = actividad.titulo.toLowerCase().includes(filtros.busqueda.toLowerCase());
            const coincideDescripcion = actividad.descripcion.toLowerCase().includes(filtros.descripcion.toLowerCase());
            const coincideDia = filtros.dia === '' || actividad.dia === filtros.dia;

            return coincideBusqueda && coincideDescripcion && coincideDia;
        });

        // Ordenar
        if (ordenamiento.campo) {
            resultado.sort((a, b) => {
                let valorA = a[ordenamiento.campo];
                let valorB = b[ordenamiento.campo];

                // Convertir a minúsculas si son strings
                if (typeof valorA === 'string') valorA = valorA.toLowerCase();
                if (typeof valorB === 'string') valorB = valorB.toLowerCase();

                if (valorA < valorB) return ordenamiento.direccion === 'asc' ? -1 : 1;
                if (valorA > valorB) return ordenamiento.direccion === 'asc' ? 1 : -1;
                return 0;
            });
        }

        return resultado;
    }, [actividades, filtros, ordenamiento]);

    // Calcular paginación
    const totalPaginas = Math.ceil(actividadesFiltradas.length / ITEMS_POR_PAGINA);
    const inicio = (paginaActual - 1) * ITEMS_POR_PAGINA;
    const actividadesPaginadas = actividadesFiltradas.slice(inicio, inicio + ITEMS_POR_PAGINA);

    // ==================== USUARIOS ====================

    // Manejar cambios en filtros de usuarios
    const handleFiltroUsuariosChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFiltrosUsuarios(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));
        setPaginaActualUsuarios(1);
    };

    // Limpiar filtros de usuarios
    const handleLimpiarFiltrosUsuarios = () => {
        setFiltrosUsuarios({
            busqueda: '',
            email: '',
            username: '',
            isAdmin: ''
        });
        setPaginaActualUsuarios(1);
    };

    // Manejar ordenamiento de usuarios
    const handleOrdenarPorUsuario = (campo) => {
        setOrdenamientoUsuarios(prev => ({
            campo,
            direccion: prev.campo === campo && prev.direccion === 'asc' ? 'desc' : 'asc'
        }));
        setPaginaActualUsuarios(1);
    };

    // Filtrar y ordenar usuarios
    const usuariosFiltrados = useMemo(() => {
        let resultado = usuarios.filter(usuario => {
            const coincideNombre = (usuario.nombre + ' ' + usuario.apellido).toLowerCase().includes(filtrosUsuarios.busqueda.toLowerCase());
            const coincideEmail = usuario.email.toLowerCase().includes(filtrosUsuarios.email.toLowerCase());
            const coincideUsername = usuario.username.toLowerCase().includes(filtrosUsuarios.username.toLowerCase());
            const coincideAdmin = filtrosUsuarios.isAdmin === '' || (filtrosUsuarios.isAdmin ? usuario.is_admin : !usuario.is_admin);

            return coincideNombre && coincideEmail && coincideUsername && coincideAdmin;
        });

        // Ordenar
        if (ordenamientoUsuarios.campo) {
            resultado.sort((a, b) => {
                let valorA = a[ordenamientoUsuarios.campo];
                let valorB = b[ordenamientoUsuarios.campo];

                if (typeof valorA === 'string') valorA = valorA.toLowerCase();
                if (typeof valorB === 'string') valorB = valorB.toLowerCase();

                if (valorA < valorB) return ordenamientoUsuarios.direccion === 'asc' ? -1 : 1;
                if (valorA > valorB) return ordenamientoUsuarios.direccion === 'asc' ? 1 : -1;
                return 0;
            });
        }

        return resultado;
    }, [usuarios, filtrosUsuarios, ordenamientoUsuarios]);

    // Calcular paginación de usuarios
    const totalPaginasUsuarios = Math.ceil(usuariosFiltrados.length / ITEMS_POR_PAGINA);
    const inicioUsuarios = (paginaActualUsuarios - 1) * ITEMS_POR_PAGINA;
    const usuariosPaginados = usuariosFiltrados.slice(inicioUsuarios, inicioUsuarios + ITEMS_POR_PAGINA);

    // Handlers de usuarios
    const handleEditarUsuario = (usuario) => {
        setUsuarioEditar(usuario);
    };

    const handleEliminarUsuario = (usuario) => {
        setUsuarioAEliminar(usuario);
    };

    const handleConfirmDeleteUsuario = async () => {
        try {
            await deleteUsuario(usuarioAEliminar.id_usuario);
            setAlertDialog({
                title: 'Usuario eliminado',
                message: 'El usuario se ha eliminado exitosamente',
                type: 'success'
            });
            setUsuarioAEliminar(null);
            setOrdenamientoUsuarios({ campo: null, direccion: 'asc' });
            setPaginaActualUsuarios(1);
        } catch (error) {
            logger.error("handleConfirmDeleteUsuario error", error);
            setAlertDialog({
                title: 'Error al eliminar',
                message: 'No se pudo eliminar el usuario. Por favor, intenta de nuevo más tarde.',
                type: 'error'
            });
            setUsuarioAEliminar(null);
        }
    };

    const handleCancelDeleteUsuario = () => {
        setUsuarioAEliminar(null);
    };

    const handleToggleAdmin = async (usuario) => {
        try {
            const nuevoEstadoAdmin = !usuario.is_admin;
            logger.info(`Toggling admin status for usuario ${usuario.id_usuario} from ${usuario.is_admin} to ${nuevoEstadoAdmin}`);

            // Usar el servicio directamente para actualizar
            await usuarioService.updateUsuario(usuario.id_usuario, {
                nombre: usuario.nombre,
                apellido: usuario.apellido,
                email: usuario.email,
                is_admin: nuevoEstadoAdmin
            });

            logger.info('Update successful, refreshing usuarios list');
            // Refrescar usuarios para sincronizar con el backend
            await fetchUsuarios();
            logger.info('Usuarios refreshed');

            setAlertDialog({
                title: 'Usuario actualizado',
                message: `${usuario.nombre} ${nuevoEstadoAdmin ? 'ahora es' : 'ya no es'} administrador`,
                type: 'success'
            });
        } catch (error) {
            logger.error("handleToggleAdmin error", error);
            setAlertDialog({
                title: 'Error al actualizar',
                message: 'No se pudo actualizar el usuario. Por favor, intenta de nuevo más tarde.',
                type: 'error'
            });
        }
    };

    const handleCloseUsuarioModal = () => {
        setUsuarioEditar(null);
        setMostrarAgregarUsuarioModal(false);
    };

    const handleSaveUsuario = async () => {
        const title = 'Usuario creado/actualizado';
        const message = mostrarAgregarUsuarioModal
            ? 'El usuario se ha creado exitosamente'
            : 'El usuario se ha actualizado exitosamente';
        setAlertDialog({ title, message, type: 'success' });
        fetchUsuarios();
        handleCloseUsuarioModal();
    };

    return (
        <div className="admin-container">
            <div className="admin-header">
                <h2>Panel de Administración</h2>
            </div>

            {/* Tabs + Botón Agregar */}
            <div className="admin-tabs-header">
                <div className="admin-tabs">
                    <button
                        className={`tab-button ${tabActiva === 'actividades' ? 'active' : ''}`}
                        onClick={() => {
                            setTabActiva('actividades');
                            setPaginaActual(1);
                        }}
                    >
                        Actividades
                    </button>
                    <button
                        className={`tab-button ${tabActiva === 'usuarios' ? 'active' : ''}`}
                        onClick={() => {
                            setTabActiva('usuarios');
                            setPaginaActualUsuarios(1);
                        }}
                    >
                        Usuarios
                    </button>
                </div>

                {tabActiva === 'actividades' && (
                    <button
                        className="btn-agregar"
                        onClick={() => setMostrarAgregarModal(true)}
                    >
                        <span>+</span>
                        Agregar Actividad
                    </button>
                )}

                {tabActiva === 'usuarios' && (
                    <button
                        className="btn-agregar"
                        onClick={() => setMostrarAgregarUsuarioModal(true)}
                    >
                        <span>+</span>
                        Agregar Usuario
                    </button>
                )}
            </div>

            {/* Contenido de Actividades */}
            {tabActiva === 'actividades' && (
                <>


                    {/* Filtros */}
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

                        {/* Controles de paginación */}
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
                </>
            )}

            {/* Contenido de Usuarios */}
            {tabActiva === 'usuarios' && (
                <>
                    {/* Filtros de Usuarios */}
                    <div className="admin-filters">
                        <div className="filter-bar-container">
                            <div className="filter-bar-header">
                                <h3 className="filter-title">Filtros</h3>
                            </div>
                            <fieldset className="filter-fieldset">
                                <legend className="sr-only">Filtrar usuarios</legend>
                                <div className="filter-inputs-row">
                                    <div className="filter-group search-group">
                                        <label htmlFor="busqueda-usuarios" className="sr-only">
                                            Buscar por nombre
                                        </label>
                                        <input
                                            type="text"
                                            id="busqueda-usuarios"
                                            name="busqueda"
                                            placeholder="Buscar por nombre..."
                                            value={filtrosUsuarios.busqueda}
                                            onChange={handleFiltroUsuariosChange}
                                            className="filter-input"
                                        />
                                    </div>

                                    {(filtrosUsuarios.busqueda || filtrosUsuarios.email || filtrosUsuarios.username || filtrosUsuarios.isAdmin !== '') && (
                                        <button
                                            onClick={handleLimpiarFiltrosUsuarios}
                                            className="filter-btn-clear"
                                            title="Limpiar"
                                        >
                                            Limpiar ✖️
                                        </button>
                                    )}

                                    <div className="filter-group">
                                        <label htmlFor="email-usuarios" className="sr-only">
                                            Email
                                        </label>
                                        <input
                                            type="text"
                                            id="email-usuarios"
                                            name="email"
                                            placeholder="Filtrar por email..."
                                            value={filtrosUsuarios.email}
                                            onChange={handleFiltroUsuariosChange}
                                            className="filter-input"
                                        />
                                    </div>

                                    <div className="filter-group">
                                        <label htmlFor="username-usuarios" className="sr-only">
                                            Username
                                        </label>
                                        <input
                                            type="text"
                                            id="username-usuarios"
                                            name="username"
                                            placeholder="Filtrar por usuario..."
                                            value={filtrosUsuarios.username}
                                            onChange={handleFiltroUsuariosChange}
                                            className="filter-input"
                                        />
                                    </div>

                                    <div className="checkbox-group">
                                        <label>
                                            <input
                                                type="checkbox"
                                                name="isAdmin"
                                                checked={filtrosUsuarios.isAdmin !== ''}
                                                onChange={(e) => {
                                                    setFiltrosUsuarios(prev => ({
                                                        ...prev,
                                                        isAdmin: e.target.checked ? true : ''
                                                    }));
                                                    setPaginaActualUsuarios(1);
                                                }}
                                            />
                                            <span>Solo administradores</span>
                                        </label>
                                    </div>
                                </div>
                            </fieldset>
                        </div>
                    </div>

                    <div className="admin-table-container">
                        <table className="admin-table">
                            <thead>
                                <tr>
                                    <th
                                        onClick={() => handleOrdenarPorUsuario('id_usuario')}
                                        className={`sortable ${ordenamientoUsuarios.campo === 'id_usuario' ? `sorted-${ordenamientoUsuarios.direccion}` : ''}`}
                                    >
                                        ID
                                    </th>
                                    <th
                                        onClick={() => handleOrdenarPorUsuario('nombre')}
                                        className={`sortable ${ordenamientoUsuarios.campo === 'nombre' ? `sorted-${ordenamientoUsuarios.direccion}` : ''}`}
                                    >
                                        Nombre
                                    </th>
                                    <th
                                        onClick={() => handleOrdenarPorUsuario('apellido')}
                                        className={`sortable ${ordenamientoUsuarios.campo === 'apellido' ? `sorted-${ordenamientoUsuarios.direccion}` : ''}`}
                                    >
                                        Apellido
                                    </th>
                                    <th
                                        onClick={() => handleOrdenarPorUsuario('username')}
                                        className={`sortable ${ordenamientoUsuarios.campo === 'username' ? `sorted-${ordenamientoUsuarios.direccion}` : ''}`}
                                    >
                                        Usuario
                                    </th>
                                    <th
                                        onClick={() => handleOrdenarPorUsuario('email')}
                                        className={`sortable ${ordenamientoUsuarios.campo === 'email' ? `sorted-${ordenamientoUsuarios.direccion}` : ''}`}
                                    >
                                        Email
                                    </th>
                                    <th
                                        onClick={() => handleOrdenarPorUsuario('is_admin')}
                                        className={`sortable ${ordenamientoUsuarios.campo === 'is_admin' ? `sorted-${ordenamientoUsuarios.direccion}` : ''}`}
                                    >
                                        Admin
                                    </th>
                                    <th className="acciones-column">Acciones</th>
                                </tr>
                            </thead>
                            <tbody>
                                {usuariosPaginados.map((usuario) => (
                                    <tr key={usuario.id_usuario}>
                                        <td>{usuario.id_usuario}</td>
                                        <td>{usuario.nombre}</td>
                                        <td>{usuario.apellido}</td>
                                        <td>{usuario.username}</td>
                                        <td>{usuario.email}</td>
                                        <td>
                                            <input
                                                type="checkbox"
                                                checked={usuario.is_admin}
                                                onChange={() => handleToggleAdmin(usuario)}
                                                className="admin-checkbox"
                                                title="Marcar/desmarcar como administrador"
                                            />
                                        </td>
                                        <td className="acciones-column">
                                            <button
                                                className="action-button edit-button"
                                                onClick={() => handleEditarUsuario(usuario)}
                                                title="Editar"
                                            >
                                                ✏️
                                            </button>
                                            <button
                                                className="action-button delete-button"
                                                onClick={() => handleEliminarUsuario(usuario)}
                                                title="Eliminar"
                                            >
                                                🗑️
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>

                        {/* Controles de paginación de usuarios */}
                        {totalPaginasUsuarios > 1 && (
                            <div className="pagination-container">
                                <span className="pagination-info">
                                    Mostrando {inicioUsuarios + 1} a {Math.min(inicioUsuarios + ITEMS_POR_PAGINA, usuariosFiltrados.length)} de {usuariosFiltrados.length} usuarios
                                </span>
                                <div className="pagination-controls">
                                    <button
                                        className="pagination-btn"
                                        onClick={() => setPaginaActualUsuarios(prev => Math.max(prev - 1, 1))}
                                        disabled={paginaActualUsuarios === 1}
                                    >
                                        ← Anterior
                                    </button>
                                    <span style={{ padding: '0.5rem 1rem', color: '#2c3e50', fontWeight: '500' }}>
                                        Página {paginaActualUsuarios} de {totalPaginasUsuarios}
                                    </span>
                                    <button
                                        className="pagination-btn"
                                        onClick={() => setPaginaActualUsuarios(prev => Math.min(prev + 1, totalPaginasUsuarios))}
                                        disabled={paginaActualUsuarios === totalPaginasUsuarios}
                                    >
                                        Siguiente →
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>
                </>
            )}

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

            {usuarioEditar && (
                <UserFormModal
                    mode="edit"
                    usuario={usuarioEditar}
                    onClose={handleCloseUsuarioModal}
                    onSave={handleSaveUsuario}
                />
            )}

            {mostrarAgregarUsuarioModal && (
                <UserFormModal
                    mode="create"
                    onClose={handleCloseUsuarioModal}
                    onSave={handleSaveUsuario}
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

            {usuarioAEliminar && (
                <ConfirmDialog
                    title="Eliminar Usuario"
                    message="¿Estás seguro de que deseas eliminar este usuario? Esta acción no se puede deshacer."
                    details={`Se eliminará: "${usuarioAEliminar.nombre} ${usuarioAEliminar.apellido}"`}
                    confirmText="Eliminar"
                    cancelText="Cancelar"
                    isDangerous={true}
                    onConfirm={handleConfirmDeleteUsuario}
                    onCancel={handleCancelDeleteUsuario}
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

export default AdminPanel; 