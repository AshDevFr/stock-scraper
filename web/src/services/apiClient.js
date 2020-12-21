import axios from 'axios';

// axios.defaults.baseURL = '/api';

export const fetchConfig = async () => {
  const response = await axios.get(`/api/config`);
  return response.data;
};
