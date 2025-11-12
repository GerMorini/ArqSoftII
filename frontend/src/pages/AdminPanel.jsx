import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AdminActivitiesTab from '../components/AdminActivitiesTab';
import AdminUsersTab from '../components/AdminUsersTab';
import AdminStatisticsTab from '../components/AdminStatisticsTab';
import useCurrentUser from '../hooks/useCurrentUser';
import '../styles/AdminPanel.css';

const AdminPanel = () => {
    const [tabActiva, setTabActiva] = useState('actividades');
    const { isAdmin } = useCurrentUser();
    const navigate = useNavigate();

    useEffect(() => {
        if (!isAdmin) {
            navigate('/');
        }
    }, [isAdmin, navigate]);

    return (
        <div className="admin-container">
            <div className="admin-header-with-tabs">
                <h2>Panel de Administración</h2>
                <div className="admin-tabs">
                    <button
                        className={`tab-button ${tabActiva === 'actividades' ? 'active' : ''}`}
                        onClick={() => setTabActiva('actividades')}
                    >
                        Actividades
                    </button>
                    <button
                        className={`tab-button ${tabActiva === 'usuarios' ? 'active' : ''}`}
                        onClick={() => setTabActiva('usuarios')}
                    >
                        Usuarios
                    </button>
                    <button
                        className={`tab-button ${tabActiva === 'estadisticas' ? 'active' : ''}`}
                        onClick={() => setTabActiva('estadisticas')}
                    >
                        Estadísticas
                    </button>
                </div>
            </div>

            {tabActiva === 'actividades' && <AdminActivitiesTab />}
            {tabActiva === 'usuarios' && <AdminUsersTab />}
            {tabActiva === 'estadisticas' && <AdminStatisticsTab />}
        </div>
    );
};

export default AdminPanel; 