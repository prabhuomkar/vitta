import axiosInstance from './axiosInstance';

export const getBudgets = (year, month) =>
  axiosInstance.get(`/budgets?year=${year}&month=${month}`);

export const createBudgets = data => axiosInstance.put('/budgets', data);
