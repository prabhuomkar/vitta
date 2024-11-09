import React, { createContext, useState, useCallback, useContext } from 'react';
import {
  fetchTransactions,
  addTransaction,
  uploadTransactions,
  editTransaction,
  removeTransaction
} from '../services/transactionsService';

export const TransactionsContext = createContext();

export const TransactionsProvider = ({ children }) => {
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [accountId, setAccountId] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');

  const getTransactions = useCallback(
    async (id, query = searchQuery) => {
      setLoading(true);
      try {
        const data = await fetchTransactions(id, query);
        setTransactions(data);
        setAccountId(id);
        setError(null);
        return { success: true };
      } catch (err) {
        setError('Failed to load transactions');
        return { success: false, error: err };
      } finally {
        setLoading(false);
      }
    },
    [searchQuery]
  );

  const updateSearchQuery = query => {
    setSearchQuery(query);
    getTransactions(accountId, query);
  };

  const createTransaction = async transactionData => {
    if (!accountId) {
      return { success: false };
    }
    try {
      const newTransaction = await addTransaction(accountId, transactionData);
      setTransactions([...transactions, newTransaction]);
      setError(null);
      return { success: true };
    } catch (err) {
      setError('Failed to create transaction');
      return { success: false, error: err };
    }
  };

  const importTransactions = async file => {
    if (!accountId) {
      return { success: false, error: new Error('Account ID is missing') };
    }

    try {
      const importedTransactions = await uploadTransactions(accountId, file);

      if (Array.isArray(importedTransactions)) {
        setTransactions([...transactions, ...importedTransactions]);
      } else {
        setTransactions([...transactions]);
      }

      setError(null);
      // eslint-disable-next-line no-console
      console.log('Imported Transactions:', importedTransactions);
      return { success: true, transactions: importedTransactions };
    } catch (err) {
      setError('Failed to import transactions');
      // eslint-disable-next-line no-console
      console.log('Error during importTransactions:', err);
      return { success: false, error: err };
    }
  };

  const updateTransaction = async (transactionId, transactionData) => {
    if (!accountId) {
      return { success: false };
    }
    try {
      await editTransaction(accountId, transactionId, transactionData);
      setTransactions(
        transactions.map(transaction =>
          transaction.id === transactionId
            ? { ...transaction, ...transactionData }
            : transaction
        )
      );
      setError(null);
      return { success: true };
    } catch (err) {
      setError('Failed to update transaction');
      return { success: false, error: err };
    }
  };

  const deleteTransaction = async transactionId => {
    if (!accountId) {
      return { success: false };
    }
    try {
      await removeTransaction(accountId, transactionId);
      setTransactions(
        transactions.filter(transaction => transaction.id !== transactionId)
      );
      setError(null);
      return { success: true };
    } catch (err) {
      setError('Failed to delete transaction');
      return { success: false, error: err };
    }
  };

  return (
    <TransactionsContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={{
        transactions,
        loading,
        error,
        getTransactions,
        createTransaction,
        importTransactions,
        updateTransaction,
        deleteTransaction,
        searchQuery,
        updateSearchQuery
      }}
    >
      {children}
    </TransactionsContext.Provider>
  );
};

export const useTransactions = () => {
  return useContext(TransactionsContext);
};
