import React from 'react';
import { Link } from 'react-router-dom';
import useCurrentUser from '../hooks/useCurrentUser';
import '../styles/Footer.css';

const Footer = () => {
    const { isLoggedIn, isAdmin } = useCurrentUser();
    const currentYear = new Date().getFullYear();

    return (
        <footer className="footer-container" role="contentinfo">
            <div className="footer-content">
                {/* Sección de enlaces */}
                <div className="footer-section">
                    <h3 className="footer-section-title">Navegación</h3>
                    <nav className="footer-nav">
                        <ul>
                            <li>
                                <Link to="/" aria-label="Ir a inicio">
                                    Inicio
                                </Link>
                            </li>
                            <li>
                                <Link to="/actividades" aria-label="Ver actividades disponibles">
                                    Actividades
                                </Link>
                            </li>
                            {isLoggedIn && !isAdmin && (
                                <li>
                                    <Link to="/mis-actividades" aria-label="Ver mis actividades">
                                        Mis Actividades
                                    </Link>
                                </li>
                            )}
                            <li>
                                <Link to="/contacto" aria-label="Ir a página de contacto">
                                    Contacto
                                </Link>
                            </li>
                        </ul>
                    </nav>
                </div>

                {/* Sección de contacto */}
                <div className="footer-section">
                    <h3 className="footer-section-title">Contacto</h3>
                    <address className="footer-contact">
                        <p>
                            <strong>Email:</strong>{' '}
                            <a
                                href="mailto:info@gympro.com"
                                aria-label="Enviar email a info@gympro.com"
                            >
                                info@gympro.com
                            </a>
                        </p>
                        <p>
                            <strong>Teléfono:</strong>{' '}
                            <a
                                href="tel:+541234567890"
                                aria-label="Llamar al +54 123 456 7890"
                            >
                                +54 123 456 7890
                            </a>
                        </p>
                    </address>
                </div>

                {/* Sección de redes sociales */}
                <div className="footer-section">
                    <h3 className="footer-section-title">Síguenos</h3>
                    <div className="footer-social">
                        <a
                            href="https://facebook.com/gympro"
                            target="_blank"
                            rel="noopener noreferrer"
                            aria-label="Visitar Facebook de GymPro"
                            title="Facebook"
                        >
                            👍 Facebook
                        </a>
                        <a
                            href="https://instagram.com/gympro"
                            target="_blank"
                            rel="noopener noreferrer"
                            aria-label="Visitar Instagram de GymPro"
                            title="Instagram"
                        >
                            📷 Instagram
                        </a>
                        <a
                            href="https://x.com/gympro"
                            target="_blank"
                            rel="noopener noreferrer"
                            aria-label="Visitar X de GymPro"
                            title="X"
                        >
                            𝕏 X
                        </a>
                    </div>
                </div>
            </div>

            {/* Línea de copyright */}
            <div className="footer-bottom">
                <p className="footer-copyright">
                    &copy; {currentYear} <strong>GymPro</strong>. Todos los derechos reservados.
                </p>
                <p className="footer-credits">
                    Diseñado y desarrollado con dedicación para mejorar tu experiencia fitness
                </p>
            </div>
        </footer>
    );
};

export default Footer;