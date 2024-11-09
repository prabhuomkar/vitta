import axiosInstance from './axiosInstance';

export const getCategories = () => axiosInstance.get('/categories');

export const createCategory = data => axiosInstance.post('/categories', data);

export const updateCatrgory = (id, data) =>
  axiosInstance.patch(`/categories/${id}`, data);

export const deleteCategory = id => axiosInstance.delete(`/categories/${id}`);
