import { useState } from "react";
import '../styles/Register.css';
import { useNavigate } from "react-router-dom";
import { storeUserSession } from "./Login";
import config from '../config/env';
import PageTransition from '../components/PageTransition';

const Register = () => {
    const [formData, setFormData] = useState({
        nombre: "",
        apellido: "",
        email: "",
        username: "",
        password: "",
        confirmPassword: ""
    });
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const navigate = useNavigate();
    const USERS_URL = config.USERS_URL;

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prevState => ({
            ...prevState,
            [name]: value
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setError("");

        if (formData.password !== formData.confirmPassword) {
            setError("Las contraseñas no coinciden");
            setIsLoading(false);
            return;
        }
        
        // regex para validación de email
        if (! /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
            setError("El email ingresado no es válido");
            setIsLoading(false);
            return;
        }

        try {
            console.log(`USERS_URL = ${USERS_URL}`);
            const response = await fetch(`${USERS_URL}/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    nombre: formData.nombre.trim(),
                    apellido: formData.apellido.trim(),
                    email: formData.email.trim(),
                    username: formData.username.trim(),
                    password: formData.password
                })
            });

            if (response.ok) {
                alert("Usuario registrado exitosamente");

                const data = await response.json();
                storeUserSession(data.access_token)

                navigate("/");
            } else {
                const errorData = await response.json();
                setError(errorData.error || "Error al registrar usuario");
            }
        } catch (error) {
            setError("Error de conexión");
            console.error("Error de conexión:", error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleBack = () => {
        navigate('/');
    };

    return (
        <PageTransition>
            <div className="register-container">
                <button onClick={handleBack} className="back-button">
                    ← Inicio
                </button>
                <form className="register-form" onSubmit={handleSubmit}>
                    <h2>Registro de Usuario</h2>

                    {error && <div className="error-message">{error}</div>}

                    <div className="input-group">
                        <input
                            type="text"
                            name="nombre"
                            placeholder="Nombre"
                            value={formData.nombre}
                            onChange={handleChange}
                            disabled={isLoading}
                            required
                        />
                    </div>

                    <div className="input-group">
                        <input
                            type="text"
                            name="apellido"
                            placeholder="Apellido"
                            value={formData.apellido}
                            onChange={handleChange}
                            disabled={isLoading}
                            required
                        />
                    </div>

                    <div className="input-group">
                        <input
                            type="text"
                            name="email"
                            placeholder="Correo electrónico"
                            value={formData.email}
                            onChange={handleChange}
                            disabled={isLoading}
                            required
                        />
                    </div>

                    <div className="input-group">
                        <input
                            type="text"
                            name="username"
                            placeholder="Nombre de usuario"
                            value={formData.username}
                            onChange={handleChange}
                            disabled={isLoading}
                            required
                        />
                    </div>

                    <div className="input-group">
                        <input
                            type="password"
                            name="password"
                            placeholder="Contraseña"
                            value={formData.password}
                            onChange={handleChange}
                            disabled={isLoading}
                            required
                        />
                    </div>

                    <div className="input-group">
                        <input
                            type="password"
                            name="confirmPassword"
                            placeholder="Confirmar Contraseña"
                            value={formData.confirmPassword}
                            onChange={handleChange}
                            disabled={isLoading}
                            required
                        />
                    </div>

                    <button type="submit" disabled={isLoading}>
                        {isLoading ? "Registrando..." : "Registrarse"}
                    </button>

                    <div className="login-link">
                        ¿Ya tienes una cuenta? <a href="/login">Iniciar Sesión</a>
                    </div>
                </form>
            </div>
        </PageTransition>
    );
};

export default Register; 