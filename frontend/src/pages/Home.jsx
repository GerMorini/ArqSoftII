import '../styles/Home.css'
import gymPortada from '../assets/login/gimnasio1.jpeg'
import { useNavigate } from 'react-router-dom'

const Home = () => {
    const navigate = useNavigate()

    return (
        <div className="home-container">
            <img
                className="img-gym"
                src={gymPortada}
                alt="Gimnasio portada de GymPro"
            />

            <div className="home-content">
                <div className="hero-section">
                    <h1 className="hero-title">
                        Transforma Tu Experiencia <span className="gradient-text">Fitness</span>
                    </h1>
                    <p className="hero-subtitle">
                        Descubre una nueva forma de entrenar con GymPro. Actividades variadas,
                        entrenadores expertos y un ambiente que te motiva cada día.
                    </p>
                    <div className="cta-buttons">
                        <button
                            className="cta-primary"
                            onClick={() => navigate('/actividades')}
                        >
                            Ver Actividades
                        </button>
                        <button
                            className="cta-secondary"
                            onClick={() => navigate('/login')}
                        >
                            Registrarse
                        </button>
                    </div>
                </div>

                <div className="stats-section">
                    <div className="stat-card">
                        <div className="stat-icon">👥</div>
                        <div className="stat-number">500+</div>
                        <div className="stat-label">Miembros Activos</div>
                    </div>
                    <div className="stat-card">
                        <div className="stat-icon">💪</div>
                        <div className="stat-number">20+</div>
                        <div className="stat-label">Actividades</div>
                    </div>
                    <div className="stat-card">
                        <div className="stat-icon">🏆</div>
                        <div className="stat-number">10+</div>
                        <div className="stat-label">Años de Experiencia</div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Home;