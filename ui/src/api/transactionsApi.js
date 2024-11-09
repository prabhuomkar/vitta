import axiosInstance from './axiosInstance';

export const getTransactions = (id, query) =>
  axiosInstance.get(`/accounts/${id}/transactions?q=${query}`);

export const importTransactions = (id, file) => {
  const formData = new FormData();
  formData.append('file', file);

  return axiosInstance.put(
    `/accounts/${id}/transactions?adapter=icici`,
    formData,
    {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    }
  );
};

export const createTransaction = (id, data) =>
  axiosInstance.post(`/accounts/${id}/transactions`, data);

export const updateTransaction = (id, tId, data) =>
  axiosInstance.patch(`/accounts/${id}/transactions/${tId}`, data);

export const deleteTransaction = (id, tId) =>
  axiosInstance.delete(`/accounts/${id}/transactions/${tId}`);
