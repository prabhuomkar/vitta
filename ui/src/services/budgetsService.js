import { getBudgets, createBudgets } from '../api/budgetsApi';

export const fetchBudgets = async (year, month) => {
  try {
    const response = await getBudgets(year, month);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error fetching budgets:', error);
    throw error;
  }
};

export const addBudget = async budgetData => {
  try {
    const response = await createBudgets(budgetData);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error creating budget:', error);
    throw error;
  }
};
