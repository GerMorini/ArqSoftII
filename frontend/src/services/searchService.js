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

      const response = await fetch(`${API_URL}/activitys?${params.toString()}`);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error searching activities:', error);
      throw error;
    }
  },

  // Get activity by ID
  getActivityById: async (id) => {
    try {
      const response = await fetch(`${API_URL}/activitys/${id}`);

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('Activity not found');
        }
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching activity:', error);
      throw error;
    }
  }
};

export default searchService;
