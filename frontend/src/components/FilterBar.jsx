import React from 'react';
import { DIAS_SEMANA } from '../constants/actividadConstants';
import '../styles/FilterBar.css';

const FilterBar = ({
    filtros,
    onFiltroChange,
    onLimpiar,
    mostrarToggle = false,
    soloInscriptoDisabled = false
}) => {
    const tieneFlltrosActivos = Object.values(filtros).some(v => v);

    return (
        <div className="filter-bar-container">
            <div className="filter-bar-header">
                <h3 className="filter-title">Filtros</h3>
            </div>

            <fieldset className="filter-fieldset">
                <legend className="sr-only">Filtrar actividades</legend>

                <div className="filter-inputs-row">
                    {/* Búsqueda */}
                    <div className="filter-group search-group">
                        <label htmlFor="busqueda" className="sr-only">
                            Buscar por título
                        </label>
                        <input
                            type="text"
                            id="busqueda"
                            name="busqueda"
                            placeholder="Buscar por título..."
                            value={filtros.busqueda}
                            onChange={onFiltroChange}
                            className="filter-input"
                            aria-label="Buscar actividades por título"
                        />
                    </div>

                    {/* Botón Limpiar Filtros */}
                    {tieneFlltrosActivos && (
                        <button
                            onClick={onLimpiar}
                            className="filter-btn-clear"
                            aria-label="Limpiar todos los filtros"
                            title="Limpiar"
                        >
                            Limpiar ✖️
                        </button>
                    )}

                    {/* Filtro de descripción */}
                    <div className="filter-group">
                        <label htmlFor="descripcion" className="sr-only">
                            Descripción
                        </label>
                        <input
                            type="text"
                            id="descripcion"
                            name="descripcion"
                            placeholder="Filtrar por descripción..."
                            value={filtros.descripcion}
                            onChange={onFiltroChange}
                            className="filter-input"
                            aria-label="Filtrar por descripción de actividad"
                        />
                    </div>

                    {/* Día */}
                    <div className="filter-group">
                        <label htmlFor="dia" className="sr-only">
                            Día de la semana
                        </label>
                        <select
                            id="dia"
                            name="dia"
                            value={filtros.dia}
                            onChange={onFiltroChange}
                            className="filter-select"
                            aria-label="Filtrar actividades por día de la semana"
                        >
                            <option value="">Día...</option>
                            {DIAS_SEMANA.map((dia) => (
                                <option key={dia.value} value={dia.value}>
                                    {dia.label}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* Checkbox Solo Inscripto */}
                    {mostrarToggle && (
                        <div className="checkbox-group">
                            <label>
                                <input
                                    type="checkbox"
                                    name="soloInscripto"
                                    checked={filtros.soloInscripto}
                                    onChange={(e) =>
                                        onFiltroChange({
                                            target: {
                                                name: 'soloInscripto',
                                                type: 'checkbox',
                                                checked: e.target.checked
                                            }
                                        })
                                    }
                                    disabled={soloInscriptoDisabled}
                                    aria-label="Mostrar solo mis actividades inscritas"
                                />
                                <span>Mis actividades</span>
                            </label>
                        </div>
                    )}
                </div>
            </fieldset>
        </div>
    );
};

export default FilterBar;
