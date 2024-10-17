import React, { createContext, useState, useEffect } from 'react';
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
    } catch (err) {
      setError(err);
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
    } catch (err) {
      setError(err);
    }
  };

  const updateAccount = async (id, accountData) => {
    try {
      const updatedAccount = await editAccount(id, accountData);
      setAccounts(prevAccounts =>
        prevAccounts.map(account =>
          account.id === id ? updatedAccount : account
        )
      );
    } catch (err) {
      setError(err);
    }
  };

  const deleteAccount = async id => {
    try {
      await removeAccount(id);
      setAccounts(prevAccounts =>
        prevAccounts.filter(account => account.id !== id)
      );
    } catch (err) {
      setError(err);
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
