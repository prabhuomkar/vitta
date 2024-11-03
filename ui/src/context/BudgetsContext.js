import React, { createContext, useContext, useState } from 'react';
import { fetchBudgets, addBudget } from '../services/budgetsService';

export const BudgetsContext = createContext();

export const BudgetsProvider = ({ children }) => {
  const [budgets, setBudgets] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const getBudgets = async (year, month) => {
    try {
      setLoading(true);
      const data = await fetchBudgets(year, month);
      setBudgets(data);
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    } finally {
      setLoading(false);
    }
  };

  const createBudget = async budgetData => {
    try {
      const newBudget = await addBudget(budgetData);
      setBudgets(prevBudgets => [...prevBudgets, newBudget]);
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  return (
    <BudgetsContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={{
        budgets,
        getBudgets,
        createBudget,
        loading,
        error
      }}
    >
      {children}
    </BudgetsContext.Provider>
  );
};

export const useBugdets = () => {
  return useContext(BudgetsContext);
};
