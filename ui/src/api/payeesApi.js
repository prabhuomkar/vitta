import axiosInstance from './axiosInstance';

export const getPayees = () => axiosInstance.get('/payees');

export const createPayee = data =>
  axiosInstance.post('/payees?updateTransactions=true', data);

export const updatePayee = (id, data) =>
  axiosInstance.patch(`/payees/${id}?updateTransactions=true`, data);

export const deletePayee = id => axiosInstance.delete(`/payees/${id}`);
