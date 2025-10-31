import React from 'react';
import { DIAS_SEMANA } from '../constants/actividadConstants';
import '../styles/FilterBar.css';

const SearchBar = ({
    filtros,
    onFiltroChange,
    onLimpiar,
    onSearch,
    mostrarToggle = false,
    soloInscriptoDisabled = false,
    isSearching = false
}) => {
    const tieneFlltrosActivos = Object.values(filtros).some(v => v);

    const handleSubmit = (e) => {
        e.preventDefault();
        onSearch();
    };

    return (
        <div className="filter-bar-container">
            <div className="filter-bar-header">
                <h3 className="filter-title">Búsqueda de Actividades</h3>
            </div>

            <form onSubmit={handleSubmit}>
                <fieldset className="filter-fieldset" disabled={isSearching}>
                    <legend className="sr-only">Buscar actividades</legend>

                    <div className="filter-inputs-row">
                        {/* Búsqueda por título */}
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

                        {/* Botón Buscar */}
                        <button
                            type="submit"
                            className="filter-btn-search"
                            aria-label="Buscar actividades"
                            disabled={isSearching}
                        >
                            {isSearching ? 'Buscando...' : 'Buscar 🔍'}
                        </button>

                        {/* Botón Limpiar Filtros */}
                        {tieneFlltrosActivos && (
                            <button
                                type="button"
                                onClick={onLimpiar}
                                className="filter-btn-clear"
                                aria-label="Limpiar todos los filtros"
                                title="Limpiar"
                                disabled={isSearching}
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
            </form>
        </div>
    );
};

export default SearchBar;
