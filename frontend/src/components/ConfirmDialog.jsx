import { useEffect, useRef } from 'react';
import '../styles/ConfirmDialog.css';

const ConfirmDialog = ({
    title,
    message,
    confirmText = 'Confirmar',
    cancelText = 'Cancelar',
    onConfirm,
    onCancel,
    isDangerous = false,
    details = null
}) => {
    const dialogRef = useRef(null);

    useEffect(() => {
        const handleEscapeKey = (event) => {
            if (event.key === 'Escape') {
                onCancel();
            }
        };

        document.addEventListener('keydown', handleEscapeKey);
        dialogRef.current?.focus();

        return () => {
            document.removeEventListener('keydown', handleEscapeKey);
        };
    }, [onCancel]);

    const handleBackdropClick = (e) => {
        if (e.target === e.currentTarget) {
            onCancel();
        }
    };

    return (
        <div className="confirm-dialog-backdrop" onClick={handleBackdropClick}>
            <div
                className={`confirm-dialog ${isDangerous ? 'dangerous' : ''}`}
                role="alertdialog"
                aria-labelledby="confirm-title"
                aria-describedby="confirm-message"
                ref={dialogRef}
                tabIndex="-1"
            >
                <h2 id="confirm-title" className="confirm-dialog-title">
                    {title}
                </h2>

                <div id="confirm-message" className="confirm-dialog-message">
                    <p>{message}</p>
                    {details && <div className="confirm-dialog-details">{details}</div>}
                </div>

                <div className="confirm-dialog-actions">
                    <button
                        onClick={onCancel}
                        className="confirm-btn-cancel"
                        aria-label={`Cancelar: ${title}`}
                    >
                        {cancelText}
                    </button>
                    <button
                        onClick={onConfirm}
                        className={`confirm-btn-confirm ${isDangerous ? 'dangerous' : ''}`}
                        aria-label={`Confirmar: ${title}`}
                    >
                        {confirmText}
                    </button>
                </div>

                <p className="confirm-dialog-hint">
                    Presiona ESC para cerrar este diálogo
                </p>
            </div>
        </div>
    );
};

export default ConfirmDialog;
