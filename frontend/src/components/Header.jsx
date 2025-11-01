import { useNavigate, useLocation } from "react-router-dom";
import { useState } from "react";
import ConfirmDialog from "./ConfirmDialog";
import useCurrentUser from "../hooks/useCurrentUser";
import "../styles/Header.css";
import { usuarioService } from "../services/usuarioService"

const Header = () => {
    const { isLoggedIn, isAdmin, username } = useCurrentUser();
    const navigate = useNavigate();
    const location = useLocation();
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
    const [showLogoutDialog, setShowLogoutDialog] = useState(false);

    const handleLogoutClick = () => {
        setShowLogoutDialog(true);
    };

    const handleConfirmLogout = () => {
        usuarioService.clearUserSession();
        setShowLogoutDialog(false);
        setIsMobileMenuOpen(false);
        navigate("/login");
    };

    const handleCancelLogout = () => {
        setShowLogoutDialog(false);
    };

    const isActive = (path) => location.pathname === path;

    const handleNavClick = (path) => {
        navigate(path);
        setIsMobileMenuOpen(false);
    };

    return (
        <header className="header-container">
            <div className="header-wrapper">
                <button
                    className="header-logo"
                    onClick={() => handleNavClick("/")}
                    aria-label="GymPro - Ir a página de inicio"
                    title="Volver a inicio"
                >
                    <span className="logo-icon">💪</span>
                    <span className="logo-text">GymPro</span>
                </button>

                {/* Hamburger menu para mobile */}
                <button
                    className="menu-toggle"
                    onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                    aria-label="Abrir menú"
                    aria-expanded={isMobileMenuOpen}
                >
                    <span className="hamburger"></span>
                </button>

                {/* Navegación */}
                <nav
                    className={`header-nav ${isMobileMenuOpen ? "open" : ""}`}
                    aria-label="Navegación principal"
                >
                    <ul className="nav-list">
                        <li>
                            <button
                                onClick={() => handleNavClick("/")}
                                className={`nav-link ${isActive("/") ? "active" : ""}`}
                                aria-label="Ir a página de inicio"
                                title="Inicio"
                            >
                                🏠 Inicio
                            </button>
                        </li>
                        <li>
                            <button
                                onClick={() => handleNavClick("/actividades")}
                                className={`nav-link ${isActive("/actividades") ? "active" : ""}`}
                                aria-label="Ver actividades disponibles"
                                title="Actividades"
                            >
                                🏋️ Actividades
                            </button>
                        </li>
                        {isLoggedIn && !isAdmin && (
                            <li>
                                <button
                                    onClick={() => handleNavClick("/mis-actividades")}
                                    className={`nav-link ${isActive("/mis-actividades") ? "active" : ""}`}
                                    aria-label="Ver mis actividades"
                                    title="Mis Actividades"
                                >
                                    📋 Mis Actividades
                                </button>
                            </li>
                        )}
                        {isAdmin && (
                            <li>
                                <button
                                    onClick={() => handleNavClick("/admin")}
                                    className={`nav-link admin-link ${isActive("/admin") ? "active" : ""}`}
                                    aria-label="Acceder al panel de administración"
                                    title="Panel Admin"
                                >
                                    ⚙️ Admin
                                </button>
                            </li>
                        )}
                    </ul>

                    {/* Sección de autenticación */}
                    <div className="auth-section">
                        {isLoggedIn ? (
                            <button
                                onClick={handleLogoutClick}
                                className="nav-link logout-btn"
                                aria-label="Cerrar sesión"
                                title="Cerrar sesión"
                            >
                                ✖️ Cerrar sesión
                            </button>
                        ) : (
                            <button
                                onClick={() => handleNavClick("/login")}
                                className={`nav-link login-btn ${isActive("/login") ? "active" : ""}`}
                                aria-label="Iniciar sesión"
                                title="Iniciar sesión"
                            >
                                🔐 Iniciar sesión
                            </button>
                        )}
                    </div>
                </nav>
            </div>

            {/* Dialog de confirmación de logout */}
            {showLogoutDialog && (
                <ConfirmDialog
                    title="Confirmar cierre de sesión"
                    message="¿Estás seguro de que deseas cerrar la sesión?"
                    details={`Usuario actual: ${username}`}
                    confirmText="Cerrar sesión"
                    cancelText="Cancelar"
                    onConfirm={handleConfirmLogout}
                    onCancel={handleCancelLogout}
                />
            )}
        </header>
    );
};

export default Header;