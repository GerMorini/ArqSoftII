import '../styles/Home.css'
// import gymPortada from '../assets/home/maxresdefault.jpg'
import gymPortada from '../assets/login/gimnasio1.jpeg'
import { useNavigate } from 'react-router-dom'

const Home = () => {
    const navigate = useNavigate();

    const handleClick = () => {
        navigate('/actividades');
    };

    return (
        <div className="home-container">
            <img 
                className="img-gym"
                src={gymPortada}
                alt="Gimnasio portada"
            />
            <button 
                className="cta-button"
                onClick={handleClick}
            >
                Actividades
            </button>
        </div>
    );
};

export default Home;