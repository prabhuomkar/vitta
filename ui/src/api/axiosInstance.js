import axios from 'axios';

const baseURL = process.env.REACT_APP_BASE_URL;
const username = process.env.REACT_APP_USERNAME;
const password = process.env.REACT_APP_PASSWORD;

// eslint-disable-next-line prefer-template
const basicAuth = 'Basic ' + btoa(`${username}:${password}`);

const axiosInstance = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
    Authorization: basicAuth
  }
});

// Request interceptor
axiosInstance.interceptors.request.use(
  config => {
    return config;
  },
  error => Promise.reject(error)
);

// Response interceptor
axiosInstance.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      // eslint-disable-next-line no-console
      console.error('Unauthorized: Please check your credentials.');
    }
    return Promise.reject(error);
  }
);

export default axiosInstance;
