import React, { useState, useEffect, useRef, useMemo } from 'react';
import '../styles/UserSelector.css';
import { usuarioService } from '../services/usuarioService';
import logger from '../utils/logger';

const UserSelector = ({ value = [], onChange, disabled = false, placeholder = "Buscar usuarios..." }) => {
    const [users, setUsers] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const dropdownRef = useRef(null);

    // Cargar usuarios al montar el componente
    useEffect(() => {
        const fetchUsers = async () => {
            try {
                setIsLoading(true);
                const usuarios = await usuarioService.getUsuarios();
                setUsers(usuarios);
                setError(null);
            } catch (err) {
                logger.error('UserSelector: Error al cargar usuarios', err);
                setError('Error al cargar usuarios');
            } finally {
                setIsLoading(false);
            }
        };

        fetchUsers();
    }, []);

    // Cerrar dropdown al hacer click fuera
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
                setIsDropdownOpen(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    // Debug: loggear cuando cambian users o value
    useEffect(() => {
        logger.info('UserSelector - value prop:', value);
        logger.info('UserSelector - users cargados:', users.length);
    }, [value, users]);

    // Obtener usuarios seleccionados (memoizado para evitar recálculos innecesarios)
    const selectedUsers = useMemo(() => {
        const result = users.filter(user => value.includes(user.id_usuario));
        logger.info('UserSelector - selectedUsers calculado:', result.length, 'de', value.length, 'IDs');
        return result;
    }, [users, value]);

    // Filtrar usuarios disponibles (excluir seleccionados)
    const availableUsers = useMemo(() =>
        users.filter(user => !value.includes(user.id_usuario)),
        [users, value]
    );

    // Aplicar búsqueda sobre usuarios disponibles
    const filteredUsers = availableUsers.filter(user => {
        const search = searchTerm.toLowerCase();
        return (
            user.username.toLowerCase().includes(search) ||
            user.nombre.toLowerCase().includes(search) ||
            user.apellido.toLowerCase().includes(search) ||
            user.id_usuario.toString().includes(search)
        );
    });

    // Formatear nombre de usuario: username (Apellido, Nombre) [id_usuario]
    const formatUserDisplay = (user) => {
        return `${user.username} (${user.apellido}, ${user.nombre}) [${user.id_usuario}]`;
    };

    // Agregar usuario a la selección
    const handleAddUser = (userId) => {
        if (!value.includes(userId)) {
            onChange([...value, userId]);
            setSearchTerm('');
            setIsDropdownOpen(false);
        }
    };

    // Remover usuario de la selección
    const handleRemoveUser = (userId) => {
        onChange(value.filter(id => id !== userId));
    };

    // Manejar cambio en el input de búsqueda
    const handleSearchChange = (e) => {
        setSearchTerm(e.target.value);
        setIsDropdownOpen(true);
    };

    // Manejar focus en el input
    const handleInputFocus = () => {
        if (!disabled) {
            setIsDropdownOpen(true);
        }
    };

    return (
        <div className="user-selector" ref={dropdownRef}>
            {/* Chips de usuarios seleccionados */}
            {selectedUsers.length > 0 && (
                <div className="selected-users-chips">
                    {selectedUsers.map(user => (
                        <div key={user.id_usuario} className="user-chip">
                            <span className="user-chip-text">{formatUserDisplay(user)}</span>
                            <button
                                type="button"
                                className="user-chip-remove"
                                onClick={() => handleRemoveUser(user.id_usuario)}
                                disabled={disabled}
                                aria-label="Remover usuario"
                            >
                                ×
                            </button>
                        </div>
                    ))}
                </div>
            )}

            {/* Input de búsqueda */}
            <div className="search-input-container">
                <input
                    type="text"
                    className="search-input"
                    value={searchTerm}
                    onChange={handleSearchChange}
                    onFocus={handleInputFocus}
                    placeholder={placeholder}
                    disabled={disabled || isLoading}
                />
                {isLoading && <span className="loading-indicator">Cargando...</span>}
                {error && <span className="error-indicator">{error}</span>}
            </div>

            {/* Dropdown de usuarios disponibles */}
            {isDropdownOpen && !disabled && !isLoading && !error && (
                <div className="users-dropdown">
                    {filteredUsers.length > 0 ? (
                        <ul className="users-list">
                            {filteredUsers.map(user => (
                                <li
                                    key={user.id_usuario}
                                    className="user-item"
                                    onClick={() => handleAddUser(user.id_usuario)}
                                >
                                    <span className="user-item-text">{formatUserDisplay(user)}</span>
                                </li>
                            ))}
                        </ul>
                    ) : (
                        <div className="no-results">
                            {searchTerm ? 'No se encontraron usuarios' : 'Todos los usuarios están seleccionados'}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

export default UserSelector;
