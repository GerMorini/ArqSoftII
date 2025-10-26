import '../styles/Home.css'
import gymPortada from '../assets/login/gimnasio1.jpeg'

const Home = () => {
    return (
        <div className="home-container">
            <img
                className="img-gym"
                src={gymPortada}
                alt="Gimnasio portada de GymPro"
            />
        </div>
    );
};

export default Home;