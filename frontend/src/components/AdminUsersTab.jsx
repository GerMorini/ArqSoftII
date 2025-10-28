import { useState, useMemo, useEffect } from 'react';
import UserFormModal from './UserFormModal';
import ConfirmDialog from './ConfirmDialog';
import AlertDialog from './AlertDialog';
import '../styles/AdminPanel.css';
import '../styles/FilterBar.css';
import { useUsuarios } from '../hooks/useUsuarios';
import { usuarioService } from '../services/usuarioService';
import logger from '../utils/logger';

const AdminUsersTab = () => {
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
    const ITEMS_POR_PAGINA = 5;
    const [alertDialog, setAlertDialog] = useState(null);

    const {
        usuarios,
        fetchUsuarios,
        deleteUsuario,
    } = useUsuarios();

    useEffect(() => {
        fetchUsuarios();
    }, [fetchUsuarios]);

    const handleFiltroUsuariosChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFiltrosUsuarios(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));
        setPaginaActualUsuarios(1);
    };

    const handleLimpiarFiltrosUsuarios = () => {
        setFiltrosUsuarios({
            busqueda: '',
            email: '',
            username: '',
            isAdmin: ''
        });
        setPaginaActualUsuarios(1);
    };

    const handleOrdenarPorUsuario = (campo) => {
        setOrdenamientoUsuarios(prev => ({
            campo,
            direccion: prev.campo === campo && prev.direccion === 'asc' ? 'desc' : 'asc'
        }));
        setPaginaActualUsuarios(1);
    };

    const usuariosFiltrados = useMemo(() => {
        let resultado = usuarios.filter(usuario => {
            const coincideNombre = (usuario.nombre + ' ' + usuario.apellido).toLowerCase().includes(filtrosUsuarios.busqueda.toLowerCase());
            const coincideEmail = usuario.email.toLowerCase().includes(filtrosUsuarios.email.toLowerCase());
            const coincideUsername = usuario.username.toLowerCase().includes(filtrosUsuarios.username.toLowerCase());
            const coincideAdmin = filtrosUsuarios.isAdmin === '' || (filtrosUsuarios.isAdmin ? usuario.is_admin : !usuario.is_admin);

            return coincideNombre && coincideEmail && coincideUsername && coincideAdmin;
        });

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

    const totalPaginasUsuarios = Math.ceil(usuariosFiltrados.length / ITEMS_POR_PAGINA);
    const inicioUsuarios = (paginaActualUsuarios - 1) * ITEMS_POR_PAGINA;
    const usuariosPaginados = usuariosFiltrados.slice(inicioUsuarios, inicioUsuarios + ITEMS_POR_PAGINA);

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

            await usuarioService.updateUsuario(usuario.id_usuario, {
                nombre: usuario.nombre,
                apellido: usuario.apellido,
                email: usuario.email,
                is_admin: nuevoEstadoAdmin
            });

            logger.info('Update successful, refreshing usuarios list');
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
        <>
            <div className="admin-tabs-header">
                <button
                    className="btn-agregar"
                    onClick={() => setMostrarAgregarUsuarioModal(true)}
                >
                    <span>+</span>
                    Agregar Usuario
                </button>
            </div>

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
        </>
    );
};

export default AdminUsersTab;
