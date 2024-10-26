import {
  getTransactions,
  createTransaction,
  importTransactions,
  updateTransaction,
  deleteTransaction
} from '../api/transactionsApi';

// Fetch all transactions for a specific account
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

// Create a new transaction for a specific account
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

// export const uploadTransactions = async (accountId, file) => {
//   try {
//     const response = await importTransactions(accountId, file);
//     return response.data;
//   } catch (error) {
//     // eslint-disable-next-line no-console
//     console.log('Error importing transaction:', error);
//     throw error;
//   }
// };

export const uploadTransactions = async (accountId, file) => {
  try {
    const response = await importTransactions(accountId, file);
    console.log('API Response:', response); // Log the entire response for debugging

    if (response.status !== 200) {
      throw new Error(`Unexpected status code: ${response.status}`);
    }

    return response.data; // Return the data from the response if status is 200
  } catch (error) {
    console.log('Error importing transaction:', error);
    throw error; // Rethrow the error to be caught in importTransactions
  }
};

// Update a transaction for a specific account
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

// Delete a transaction for a specific account
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
