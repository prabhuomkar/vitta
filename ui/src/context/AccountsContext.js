import React, { createContext, useContext, useState, useEffect } from 'react';
import {
  fetchAccounts,
  addAccount,
  editAccount,
  removeAccount
} from '../services/accountsService';

export const AccountsContext = createContext();

export const AccountsProvider = ({ children }) => {
  const [accounts, setAccounts] = useState([]);
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

  useEffect(() => {
    getAccounts();
  }, []);

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

  return (
    <AccountsContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={{
        accounts,
        getAccounts,
        loading,
        error,
        createAccount,
        updateAccount,
        deleteAccount
      }}
    >
      {children}
    </AccountsContext.Provider>
  );
};

export const useAccounts = () => {
  return useContext(AccountsContext);
};
