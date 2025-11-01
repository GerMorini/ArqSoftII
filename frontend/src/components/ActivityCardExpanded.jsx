import '../styles/ActivityCardExpanded.css';
import { useEscapeKey } from '../hooks/useEscapeKey';
import logoGym from '../../img/icon-gym2.png';

const ActivityCardExpanded = ({ actividad, onClose }) => {
    useEscapeKey(onClose);

    const handleBackdropClick = (e) => {
        if (e.target === e.currentTarget) {
            onClose();
        }
    };

    return (
        <>
            <div className="activity-expanded-overlay" onClick={handleBackdropClick} />
            <div className="activity-expanded-card">
                <button className="activity-expanded-close" onClick={onClose}>
                    ✕
                </button>

                <h2 className="activity-expanded-title">{actividad.titulo}</h2>

                <div className="activity-expanded-content">
                    <div className="activity-expanded-image">
                        <img
                            src={actividad.foto_url || logoGym}
                            alt={actividad.titulo}
                            onError={(e) => { e.target.src = logoGym; }}
                        />
                    </div>

                    <div className="activity-expanded-details">
                        <div className="activity-expanded-section">
                            <h3>Descripción</h3>
                            <p>{actividad.descripcion}</p>
                        </div>

                        <div className="activity-expanded-section">
                            <h3>Información</h3>
                            <p><strong>Instructor:</strong> {actividad.instructor}</p>
                            <p><strong>Día:</strong> {actividad.dia}</p>
                            <p><strong>Horario:</strong> {actividad.hora_inicio} a {actividad.hora_fin}</p>
                            <p><strong>Lugares disponibles:</strong> {actividad.lugares_disponibles}</p>
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
};

export default ActivityCardExpanded;
