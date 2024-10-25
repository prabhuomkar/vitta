import {
  getPayees,
  createPayee,
  updatePayee,
  deletePayee
} from '../api/payeesApi';

export const fetchPayees = async () => {
  try {
    const response = await getPayees();
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error fetching payees:', error);
    throw error;
  }
};

export const addPayee = async payeeData => {
  try {
    const response = await createPayee(payeeData);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error creating payee:', error);
    throw error;
  }
};

export const editPayee = async (id, payeeData) => {
  try {
    const response = await updatePayee(id, payeeData);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error updating payee:', error);
    throw error;
  }
};

export const removePayee = async id => {
  try {
    const response = await deletePayee(id);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error deleting payee:', error);
    throw error;
  }
};
