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
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [accountId, setAccountId] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [hasNextPage, setHasNextPage] = useState(false);
  const totalPages = Math.ceil(total / limit);

  const getTransactions = useCallback(
    async (id, query = searchQuery, pageNum = page, pageLimit = limit) => {
      setLoading(true);
      try {
        const data = await fetchTransactions(id, {
          query,
          page: pageNum,
          limit: pageLimit
        });
        setTransactions(data.transactions || []);
        setTotal(data.total || 0);
        setAccountId(id);
        setError(null);

        const totalPageCount = Math.ceil(data.total / pageLimit);
        setHasNextPage(pageNum < totalPageCount);

        return { success: true };
      } catch (err) {
        setError('Failed to load transactions');
        return { success: false, error: err };
      } finally {
        setLoading(false);
      }
    },
    [searchQuery, page, limit]
  );

  const updateSearchQuery = query => {
    setSearchQuery(query);
    setPage(1);
  };

  const goToNextPage = () => {
    if (hasNextPage) setPage(prev => prev + 1);
  };

  const goToPreviousPage = () => {
    setPage(prev => (prev > 1 ? prev - 1 : 1));
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

  const importTransactions = async (accId, file) => {
    if (!accId) {
      return { success: false, error: new Error('Account ID is missing') };
    }

    try {
      const response = await uploadTransactions(accId, file);
      const { importedTransactions, status } = response;

      if (status === 200 && Array.isArray(importedTransactions)) {
        setTransactions([...transactions, ...importedTransactions]);
      } else if (status === 200) {
        setTransactions([...transactions]);
      }

      setError(null);
      // eslint-disable-next-line no-console
      return { success: true, status, transactions: importedTransactions };
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
        total,
        totalPages,
        loading,
        error,
        getTransactions,
        createTransaction,
        importTransactions,
        updateTransaction,
        deleteTransaction,
        searchQuery,
        updateSearchQuery,
        goToNextPage,
        goToPreviousPage,
        page,
        setPage,
        hasNextPage
      }}
    >
      {children}
    </TransactionsContext.Provider>
  );
};

export const useTransactions = () => {
  return useContext(TransactionsContext);
};
