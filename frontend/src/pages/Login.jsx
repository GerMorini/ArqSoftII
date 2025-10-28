import { useState } from "react";
import '../styles/Login.css';
import { useNavigate } from "react-router-dom";
import PageTransition from '../components/PageTransition';
import AlertDialog from '../components/AlertDialog';
import { useUsuarios } from '../hooks/useUsuarios';

const Login = () => {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [alertDialog, setAlertDialog] = useState(null);
    const navigate = useNavigate();
    const { login, loading } = useUsuarios();

    const handleLogin = async (e) => {
        e.preventDefault();

        try {
            await login(username, password);
            navigate("/actividades");
        } catch (err) {
            setAlertDialog({
                title: "Error de autenticación",
                message: err.message || "Usuario o contraseña incorrectos",
                type: "error"
            });
        }
    };

    const handleAlertClose = () => {
        setAlertDialog(null);
    };

    const handleBack = () => {
        navigate('/');
    };

    return (
        <PageTransition>
            <div className="login-container">
                <button onClick={handleBack} className="back-button">
                    ← Inicio
                </button>
                <form className="login-form" onSubmit={handleLogin}>
                    <h2>Iniciar Sesión</h2>

                    <div className="input-group">
                        <input
                            type="text"
                            placeholder="Usuario"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            disabled={loading}
                            required
                        />
                    </div>

                    <div className="input-group">
                        <input
                            type="password"
                            placeholder="Contraseña"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            disabled={loading}
                            required
                        />
                    </div>

                    <button type="submit" disabled={loading}>
                        {loading ? "Ingresando..." : "Ingresar"}
                    </button>

                    <div className="register-link">
                        ¿No tienes una cuenta? <a href="/register">Regístrate ahora</a>
                    </div>
                </form>
            </div>

            {alertDialog && (
                <AlertDialog
                    title={alertDialog.title}
                    message={alertDialog.message}
                    type={alertDialog.type}
                    onClose={handleAlertClose}
                />
            )}
        </PageTransition>
    );
};

export default Login;
