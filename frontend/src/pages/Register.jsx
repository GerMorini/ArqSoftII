import { useState } from "react";
import '../styles/Register.css';
import { useNavigate } from "react-router-dom";
import PageTransition from '../components/PageTransition';
import AlertDialog from '../components/AlertDialog';
import { useUsuarios } from '../hooks/useUsuarios';
import { usuarioService } from '../services/usuarioService';

const Register = () => {
    const [formData, setFormData] = useState({
        nombre: "",
        apellido: "",
        email: "",
        username: "",
        password: "",
        confirmPassword: ""
    });
    const [alertDialog, setAlertDialog] = useState(null);
    const [validationErrors, setValidationErrors] = useState({});
    const navigate = useNavigate();
    const { register, loading } = useUsuarios();

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prevState => ({
            ...prevState,
            [name]: value
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setValidationErrors({});

        try {
            // Validar formulario
            const errors = usuarioService.validateUsuarioForm(formData, true);

            // Validar que las contraseñas coincidan
            if (formData.password !== formData.confirmPassword) {
                errors.confirmPassword = 'Las contraseñas no coinciden';
            }

            if (Object.keys(errors).length > 0) {
                setValidationErrors(errors);
                return;
            }

            const { confirmPassword, ...registerData } = formData;
            await register(registerData);
            navigate("/");
        } catch (err) {
            setAlertDialog({
                title: "Error al registrar",
                message: err.message || "Error al registrar usuario",
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
            <div className="register-container">
                <button onClick={handleBack} className="back-button">
                    ← Inicio
                </button>
                <form className="register-form" onSubmit={handleSubmit}>
                    <h2>Registro de Usuario</h2>

                    <div className="input-group">
                        <input
                            type="text"
                            name="nombre"
                            placeholder="Nombre"
                            value={formData.nombre}
                            onChange={handleChange}
                            disabled={loading}
                            required
                        />
                        {validationErrors.nombre && <span className="error-text">{validationErrors.nombre}</span>}
                    </div>

                    <div className="input-group">
                        <input
                            type="text"
                            name="apellido"
                            placeholder="Apellido"
                            value={formData.apellido}
                            onChange={handleChange}
                            disabled={loading}
                            required
                        />
                        {validationErrors.apellido && <span className="error-text">{validationErrors.apellido}</span>}
                    </div>

                    <div className="input-group">
                        <input
                            type="text"
                            name="email"
                            placeholder="Correo electrónico"
                            value={formData.email}
                            onChange={handleChange}
                            disabled={loading}
                            required
                        />
                        {validationErrors.email && <span className="error-text">{validationErrors.email}</span>}
                    </div>

                    <div className="input-group">
                        <input
                            type="text"
                            name="username"
                            placeholder="Nombre de usuario"
                            value={formData.username}
                            onChange={handleChange}
                            disabled={loading}
                            required
                        />
                        {validationErrors.username && <span className="error-text">{validationErrors.username}</span>}
                    </div>

                    <div className="input-group">
                        <input
                            type="password"
                            name="password"
                            placeholder="Contraseña"
                            value={formData.password}
                            onChange={handleChange}
                            disabled={loading}
                            required
                        />
                        {validationErrors.password && <span className="error-text">{validationErrors.password}</span>}
                    </div>

                    <div className="input-group">
                        <input
                            type="password"
                            name="confirmPassword"
                            placeholder="Confirmar Contraseña"
                            value={formData.confirmPassword}
                            onChange={handleChange}
                            disabled={loading}
                            required
                        />
                        {validationErrors.confirmPassword && <span className="error-text">{validationErrors.confirmPassword}</span>}
                    </div>

                    <button type="submit" disabled={loading}>
                        {loading ? "Registrando..." : "Registrarse"}
                    </button>

                    <div className="login-link">
                        ¿Ya tienes una cuenta? <a href="/login">Iniciar Sesión</a>
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

export default Register;
