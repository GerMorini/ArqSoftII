import { useNavigate } from "react-router-dom";
import "../styles/Header.css";


const Header = ( ) => {
    const isLoggedIn = localStorage.getItem("isLoggedIn") === "true";
    const isAdmin = localStorage.getItem("isAdmin") === "true";
    const navigate = useNavigate();
    const logout = () => {
        localStorage.removeItem("isLoggedIn");
        localStorage.removeItem("isAdmin");
        localStorage.removeItem("access_token");
        navigate("/");
    }

    return (
        <header>
            <div className="header-container">
                <nav className="header-content">
                    <h1 className="header-title" onClick={() => navigate("/")} style={{ cursor: "pointer" }}>💪 GymPro</h1> 
                    <div className="header-links"> 
                        <a href="/">Inicio 🏠</a>
                        <a href="/actividades">Actividades 🏋🏽‍♂️</a>
                        {isAdmin && (
                            <a href="/admin">Panel Admin 👨🏼‍🔧</a>
                        )}
                        {isLoggedIn ? (
                            <button onClick={logout}>Cerrar sesión ✖️</button>
                        ) : (
                            <a href="/login">Iniciar Sesión ⏫</a>
                        )}
                    </div>
                </nav>
            </div>
        </header>
    );
}

export default Header;