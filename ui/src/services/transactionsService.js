import {
  getTransactions,
  createTransaction,
  importTransactions,
  updateTransaction,
  deleteTransaction
} from '../api/transactionsApi';

export const fetchTransactions = async accountId => {
  try {
    const response = await getTransactions(accountId);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error fetching transactions:', error);
    throw error;
  }
};

export const addTransaction = async (accountId, transactionData) => {
  try {
    const response = await createTransaction(accountId, transactionData);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error creating transaction:', error);
    throw error;
  }
};

export const uploadTransactions = async (accountId, file) => {
  // eslint-disable-next-line no-useless-catch
  try {
    const response = await importTransactions(accountId, file);
    // console.log('API Response:', response);

    if (response.status !== 200) {
      throw new Error(`Unexpected status code: ${response.status}`);
    }

    return response.data;
  } catch (error) {
    // console.log('Error importing transaction:', error);
    throw error;
  }
};

export const editTransaction = async (
  accountId,
  transactionId,
  transactionData
) => {
  try {
    const response = await updateTransaction(
      accountId,
      transactionId,
      transactionData
    );
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error updating transaction:', error);
    throw error;
  }
};

export const removeTransaction = async (accountId, transactionId) => {
  try {
    const response = await deleteTransaction(accountId, transactionId);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error deleting transaction:', error);
    throw error;
  }
};
