import { useEffect, useState } from 'react';
import { getStatistics } from '../services/actividadService';
import '../styles/AdminStatisticsTab.css';

export default function AdminStatisticsTab() {
    const [statistics, setStatistics] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetchStatistics();
    }, []);

    const fetchStatistics = async () => {
        try {
            setLoading(true);
            setError(null);
            const data = await getStatistics();
            setStatistics(data);
        } catch (err) {
            console.error('Error fetching statistics:', err);
            setError('Error al cargar las estadísticas');
        } finally {
            setLoading(false);
        }
    };

    if (loading) {
        return <div className="loading">Cargando estadísticas...</div>;
    }

    if (error) {
        return <div className="error">{error}</div>;
    }

    if (!statistics) {
        return <div className="no-data">No hay datos disponibles</div>;
    }

    return (
        <div className="statistics-container">
            <h2>Estadísticas de Actividades</h2>

            <div className="stats-grid">
                {/* Card 1: Total Activities */}
                <div className="stat-card">
                    <div className="stat-icon">🎯</div>
                    <div className="stat-content">
                        <div className="stat-value">{statistics.total_actividades}</div>
                        <div className="stat-label">Actividades Totales</div>
                    </div>
                </div>

                {/* Card 2: Total Enrollments */}
                <div className="stat-card">
                    <div className="stat-icon">👥</div>
                    <div className="stat-content">
                        <div className="stat-value">{statistics.total_inscripciones}</div>
                        <div className="stat-label">Inscripciones Totales</div>
                    </div>
                </div>

                {/* Card 3: Capacity Utilization */}
                <div className="stat-card">
                    <div className="stat-icon">📊</div>
                    <div className="stat-content">
                        <div className="stat-value">{statistics.utilizacion_capacidad.toFixed(1)}%</div>
                        <div className="stat-label">Utilización de Capacidad</div>
                    </div>
                </div>

                {/* Card 4: Average Enrollment Rate */}
                <div className="stat-card">
                    <div className="stat-icon">📈</div>
                    <div className="stat-content">
                        <div className="stat-value">{statistics.tasa_promedio_inscripcion.toFixed(1)}</div>
                        <div className="stat-label">Promedio de Inscritos por Actividad</div>
                    </div>
                </div>

                {/* Card 5: Full Activities */}
                <div className="stat-card">
                    <div className="stat-icon">🔴</div>
                    <div className="stat-content">
                        <div className="stat-value">{statistics.actividades_llenas}</div>
                        <div className="stat-label">Actividades Llenas</div>
                    </div>
                </div>

                {/* Card 6: Available Activities */}
                <div className="stat-card">
                    <div className="stat-icon">🟢</div>
                    <div className="stat-content">
                        <div className="stat-value">{statistics.actividades_disponibles}</div>
                        <div className="stat-label">Actividades con Disponibilidad</div>
                    </div>
                </div>
            </div>

            {/* Most Popular Activity */}
            {statistics.actividad_mas_popular && (
                <div className="popular-activity-section">
                    <h3>Actividad Más Popular</h3>
                    <div className="popular-activity-card">
                        <div className="popular-activity-header">
                            <h4>{statistics.actividad_mas_popular.titulo}</h4>
                            <span className="popular-badge">⭐ Más Popular</span>
                        </div>
                        <p className="popular-activity-description">
                            {statistics.actividad_mas_popular.descripcion}
                        </p>
                        <div className="popular-activity-details">
                            <span>📅 {statistics.actividad_mas_popular.dia}</span>
                            <span>🕒 {statistics.actividad_mas_popular.hora_inicio} - {statistics.actividad_mas_popular.hora_fin}</span>
                            <span>👨‍🏫 {statistics.actividad_mas_popular.instructor}</span>
                            <span>
                                👥 {statistics.actividad_mas_popular.lugares_disponibles < statistics.actividad_mas_popular.cupo
                                    ? `${statistics.actividad_mas_popular.cupo - statistics.actividad_mas_popular.lugares_disponibles}/${statistics.actividad_mas_popular.cupo}`
                                    : `0/${statistics.actividad_mas_popular.cupo}`} inscritos
                            </span>
                        </div>
                    </div>
                </div>
            )}

            {/* Activities by Day of Week */}
            <div className="day-distribution-section">
                <h3>Distribución por Día de la Semana</h3>
                <div className="day-distribution-grid">
                    {statistics.actividades_por_dia && statistics.actividades_por_dia.length > 0 ? (
                        statistics.actividades_por_dia
                            .sort((a, b) => {
                                const daysOrder = ['Lunes', 'Martes', 'Miércoles', 'Jueves', 'Viernes', 'Sábado', 'Domingo'];
                                return daysOrder.indexOf(a.dia) - daysOrder.indexOf(b.dia);
                            })
                            .map(day => (
                                <div key={day.dia} className="day-bar">
                                    <div className="day-label">{day.dia}</div>
                                    <div className="day-bar-container">
                                        <div
                                            className="day-bar-fill"
                                            style={{
                                                width: `${(day.count / statistics.total_actividades) * 100}%`
                                            }}
                                        ></div>
                                        <span className="day-count">{day.count}</span>
                                    </div>
                                </div>
                            ))
                    ) : (
                        <div className="no-data">No hay datos de distribución por día</div>
                    )}
                </div>
            </div>

            {/* Refresh Button */}
            <div className="refresh-section">
                <button onClick={fetchStatistics} className="refresh-button">
                    🔄 Actualizar Estadísticas
                </button>
            </div>
        </div>
    );
}
