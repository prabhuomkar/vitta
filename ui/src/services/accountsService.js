import {
  getAccounts,
  createAccount,
  updateAccount,
  deleteAccount
} from '../api/accountsApi';

export const fetchAccounts = async () => {
  try {
    const response = await getAccounts();
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error fetching accounts:', error);
    throw error;
  }
};

export const addAccount = async accountData => {
  try {
    const response = await createAccount(accountData);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error creating account:', error);
    throw error;
  }
};

export const editAccount = async (id, accountData) => {
  try {
    const response = await updateAccount(id, accountData);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error updating account:', error);
    throw error;
  }
};

export const removeAccount = async id => {
  try {
    const response = await deleteAccount(id);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error deleting account:', error);
    throw error;
  }
};
