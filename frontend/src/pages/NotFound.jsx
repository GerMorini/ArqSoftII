import '../styles/NotFound.css'
import { useNavigate } from 'react-router-dom'
import gymPortada from '../assets/login/gimnasio1.jpeg'
import { useScrollToTop } from '../hooks/useScrollToTop'

const NotFound = () => {
    const navigate = useNavigate()
    useScrollToTop()

    return (
        <div className="notfound-container">
            <img
                className="notfound-bg-img"
                src={gymPortada}
                alt="Fondo gimnasio"
            />
            <div className="notfound-content">
                <div className="notfound-icon">🏋️</div>
                <h1 className="notfound-title">
                    <span className="error-code">404</span>
                </h1>
                <h2 className="notfound-subtitle">
                    ¡Página No Encontrada!
                </h2>
                <p className="notfound-message">
                    Parece que te has perdido en el gimnasio. Esta página no existe
                    o ha sido movida a otra ubicación.
                </p>
                <div className="notfound-buttons">
                    <button
                        className="btn-home"
                        onClick={() => navigate('/')}
                    >
                        🏠 Volver al Inicio
                    </button>
                    <button
                        className="btn-activities"
                        onClick={() => navigate('/actividades')}
                    >
                        💪 Ver Actividades
                    </button>
                </div>
            </div>
        </div>
    );
};

export default NotFound;
