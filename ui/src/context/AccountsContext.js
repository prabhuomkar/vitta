import React, { createContext, useContext, useState, useEffect } from 'react';
import {
  fetchAccounts,
  fetchAccountById,
  addAccount,
  editAccount,
  removeAccount,
  fetchAdapters
} from '../services/accountsService';

export const AccountsContext = createContext();

export const AccountsProvider = ({ children }) => {
  const [accounts, setAccounts] = useState([]);
  const [adapters, setAdapters] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const getAccounts = async () => {
    try {
      setLoading(true);
      const data = await fetchAccounts();
      setAccounts(data);
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    } finally {
      setLoading(false);
    }
  };

  const getAccountById = async id => {
    try {
      setLoading(true);
      const account = await fetchAccountById(id);
      return { success: true, account };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    } finally {
      setLoading(false);
    }
  };

  const createAccount = async accountData => {
    try {
      const newAccount = await addAccount(accountData);
      setAccounts(prevAccounts => [...prevAccounts, newAccount]);
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  const updateAccount = async (id, accountData) => {
    try {
      await editAccount(id, accountData);

      // Update the local state with the latest account data
      setAccounts(prevAccounts =>
        prevAccounts.map(account =>
          account.id === id ? { ...account, ...accountData } : account
        )
      );
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  const deleteAccount = async id => {
    try {
      await removeAccount(id);
      setAccounts(prevAccounts =>
        prevAccounts.filter(account => account.id !== id)
      );
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  const getAdapters = async () => {
    try {
      setLoading(true);
      const data = await fetchAdapters();
      setAdapters(data);
      return { success: true, adapters };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getAccounts();
  }, []);

  useEffect(() => {
    getAdapters();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <AccountsContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={{
        accounts,
        getAccounts,
        getAccountById,
        loading,
        error,
        createAccount,
        updateAccount,
        deleteAccount,
        adapters,
        getAdapters
      }}
    >
      {children}
    </AccountsContext.Provider>
  );
};

export const useAccounts = () => {
  return useContext(AccountsContext);
};
