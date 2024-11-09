import {
  getGroups,
  createGroup,
  updateGroup,
  deleteGroup
} from '../api/groupsApi';

export const fetchGroups = async () => {
  try {
    const response = await getGroups();
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error fetching groups:', error);
    throw error;
  }
};

export const addGroup = async groupData => {
  try {
    const response = await createGroup(groupData);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error creating group:', error);
    throw error;
  }
};

export const editGroup = async (id, groupData) => {
  try {
    const response = await updateGroup(id, groupData);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error updating group:', error);
    throw error;
  }
};

export const removeGroup = async id => {
  try {
    const response = await deleteGroup(id);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error deleting group:', error);
    throw error;
  }
};
