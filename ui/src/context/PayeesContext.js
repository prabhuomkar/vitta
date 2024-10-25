import React, { createContext, useContext, useState, useEffect } from 'react';
import {
  fetchPayees,
  addPayee,
  editPayee,
  removePayee
} from '../services/payeesService';

export const PayeesContext = createContext();

export const PayeesProvider = ({ children }) => {
  const [payees, setPayees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const getPayees = async () => {
    try {
      setLoading(true);
      const data = await fetchPayees();
      setPayees(data);
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getPayees();
  }, []);

  const createPayee = async payeeData => {
    try {
      const newPayee = await addPayee(payeeData);
      setPayees(prevPayees => [...prevPayees, newPayee]);
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  const updatePayee = async (id, payeeData) => {
    try {
      await editPayee(id, payeeData);

      // Update the local state with the latest payee data
      setPayees(prevPayees =>
        prevPayees.map(payee =>
          payee.id === id ? { ...payee, ...payeeData } : payee
        )
      );
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  const deletePayee = async id => {
    try {
      await removePayee(id);
      setPayees(prevPayees => prevPayees.filter(payee => payee.id !== id));
      return { success: true, id };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  return (
    <PayeesContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={{
        payees,
        getPayees,
        loading,
        error,
        createPayee,
        updatePayee,
        deletePayee
      }}
    >
      {children}
    </PayeesContext.Provider>
  );
};

export const usePayees = () => {
  return useContext(PayeesContext);
};
