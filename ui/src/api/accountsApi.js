import axiosInstance from './axiosInstance';

export const getAccounts = () => axiosInstance.get('/accounts');

export const createAccount = data => axiosInstance.post('/accounts', data);

export const updateAccount = (id, data) =>
  axiosInstance.patch(`/accounts/${id}`, data);

export const deleteAccount = id => axiosInstance.delete(`/accounts/${id}`);
