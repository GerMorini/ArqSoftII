import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import EditarActividadModal from '../components/EditarActividadModal';
import AgregarActividadModal from '../components/AgregarActividadModal';
import '../styles/AdminPanel.css';
import config from '../config/env';

const AdminPanel = () => {
    const [actividades, setActividades] = useState([]);
    const [actividadEditar, setActividadEditar] = useState(null);
    const [mostrarAgregarModal, setMostrarAgregarModal] = useState(false);
    const navigate = useNavigate();
    const ACTIVITIES_URL = config.ACTIVITIES_URL;

    useEffect(() => {
        const isAdmin = localStorage.getItem("isAdmin") === "true";
        if (!isAdmin) {
            navigate('/');
            return;
        }
        fetchActividades();
    }, [navigate]);

    const fetchActividades = async () => {
        try {
            console.log(`ACTIVITIES_URL = ${ACTIVITIES_URL}`);
            const response = await fetch(`${ACTIVITIES_URL}/activities`);
            if (response.ok) {
                const data = await response.json();
                setActividades(data.activities);
            }
        } catch (error) {
            console.error("Error al cargar actividades:", error);
        }
    };

    const handleEditar = (actividad) => {
        setActividadEditar(actividad);
    };

    const handleCloseModal = () => {
        setActividadEditar(null);
        setMostrarAgregarModal(false);
    };

    const handleSaveEdit = () => {
        fetchActividades();
        handleCloseModal();
    };

    const handleEliminar = async (actividad) => {
        if (!actividad.id_actividad) {
            console.error("Error: La actividad no tiene ID", actividad);
            alert('Error: No se puede eliminar la actividad porque no tiene ID');
            return;
        }

        if (window.confirm('¿Estás seguro de que deseas eliminar esta actividad? Se eliminarán también todas las inscripciones asociadas.')) {
            try {
                console.log(`ACTIVITIES_URL = ${ACTIVITIES_URL}`);
                const response = await fetch(`${ACTIVITIES_URL}/activities/${actividad.id_actividad}`, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
                        'Content-Type': 'application/json'
                    }
                });

                if (response.ok) {
                    fetchActividades();
                    alert('Actividad eliminada con éxito');
                } else {
                    const errorData = await response.json();
                    alert(errorData.error || 'Error al eliminar la actividad');
                }
            } catch (error) {
                console.error("Error al eliminar:", error);
                alert('Error al eliminar la actividad. Por favor, intenta de nuevo más tarde.');
            }
        }
    };

    return (
        <div className="admin-container">
            <div className="admin-header">
                <h2>Panel de Administración</h2>
                <button 
                    className="btn-agregar"
                    onClick={() => setMostrarAgregarModal(true)}
                >
                    <span>+</span>
                    Agregar Actividad
                </button>
            </div>

            <div className="admin-table-container">
                <table className="admin-table">
                    <thead>
                        <tr>
                            <th>Título</th>
                            <th>Descripción</th>
                            <th>Instructor</th>
                            <th>Día</th>
                            <th>Horario</th>
                            <th>Cupo</th>
                            <th>Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {actividades.map((actividad) => (
                            <tr key={actividad.id_actividad}>
                                <td>{actividad.titulo}</td>
                                <td>{actividad.descripcion}</td>
                                <td>{actividad.instructor}</td>
                                <td>{actividad.dia}</td>
                                <td>{actividad.hora_inicio} - {actividad.hora_fin}</td>
                                <td>{actividad.cupo - actividad.lugares} / {actividad.cupo}</td>
                                <td className="acciones-column">
                                    <button 
                                        className="action-button edit-button"
                                        onClick={() => handleEditar(actividad)}
                                        title="Editar"
                                    >
                                        ✏️
                                    </button>
                                    <button 
                                        className="action-button delete-button"
                                        onClick={() => handleEliminar(actividad)}
                                        title="Eliminar"
                                    >
                                        🗑️
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            {actividadEditar && (
                <EditarActividadModal
                    actividad={actividadEditar}
                    onClose={handleCloseModal}
                    onSave={handleSaveEdit}
                />
            )}

            {mostrarAgregarModal && (
                <AgregarActividadModal
                    onClose={handleCloseModal}
                    onSave={handleSaveEdit}
                />
            )}
        </div>
    );
};

export default AdminPanel; 