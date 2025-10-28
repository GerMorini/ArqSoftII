import '../styles/AlertDialog.css';
import { useEscapeKey } from '../hooks/useEscapeKey';

const AlertDialog = ({
    title,
    message,
    onClose,
    type = 'success'
}) => {
    useEscapeKey(onClose);

    const handleBackdropClick = (e) => {
        if (e.target === e.currentTarget) {
            onClose();
        }
    };

    return (
        <div className="alert-dialog-backdrop" onClick={handleBackdropClick}>
            <div
                className={`alert-dialog ${type}`}
                role="alert"
                aria-labelledby="alert-title"
                aria-describedby="alert-message"
                tabIndex="-1"
            >
                <div className="alert-icon">
                    {type === 'success' && '✓'}
                    {type === 'error' && '✗'}
                    {type === 'info' && 'ℹ'}
                </div>

                <h2 id="alert-title" className="alert-dialog-title">
                    {title}
                </h2>

                <div id="alert-message" className="alert-dialog-message">
                    {message}
                </div>

                <button
                    onClick={onClose}
                    className={`alert-btn ${type}`}
                    aria-label={`Cerrar: ${title}`}
                >
                    Aceptar
                </button>

                <p className="alert-dialog-hint">
                    Presiona ESC para cerrar
                </p>
            </div>
        </div>
    );
};

export default AlertDialog;
