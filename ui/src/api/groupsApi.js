import axiosInstance from './axiosInstance';

export const getGroups = () => axiosInstance.get('/groups');

export const createGroup = data => axiosInstance.post('/groups', data);

export const updateGroup = (id, data) =>
  axiosInstance.patch(`/groups/${id}`, data);

export const deleteGroup = id => axiosInstance.delete(`/groups/${id}`);
