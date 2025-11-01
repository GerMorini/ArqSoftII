import config from '../config/env';

const API_URL = config.SEARCH_URL;

export const searchService = {
  // Search activities with filters
  searchActivities: async (filters = {}) => {
    try {
      const params = new URLSearchParams();

      if (filters.titulo) params.append('titulo', filters.titulo);
      if (filters.descripcion) params.append('descripcion', filters.descripcion);
      if (filters.dia) params.append('diaSemana', filters.dia);
      if (filters.page) params.append('page', filters.page);
      if (filters.count) params.append('count', filters.count);

      const response = await fetch(`${API_URL}/activities?${params.toString()}`);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      let data = await response.json();

      // Map 'id' to 'id_actividad' for each result to match frontend expectations
      if (data.results && Array.isArray(data.results)) {
        data.results = data.results.map(activity => ({
          ...activity,
          id_actividad: activity.id
        }));
      }

      return data;
    } catch (error) {
      console.error('Error searching activities:', error);
      throw error;
    }
  },
};

export default searchService;
