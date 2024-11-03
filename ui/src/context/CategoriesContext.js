import React, { createContext, useContext, useEffect, useState } from 'react';
import {
  fetchCategories,
  addCategory,
  editCategory,
  removeCategory
} from '../services/categoriesService';

export const CategoriesContext = createContext();

export const CategoriesProvider = ({ children }) => {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const getCategories = async () => {
    try {
      setLoading(true);
      const data = await fetchCategories();
      setCategories(data);
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    getCategories();
  }, []);

  const createCategory = async categoryData => {
    try {
      const newCategory = await addCategory(categoryData);
      setCategories(prevCategories => [...prevCategories, newCategory]);
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  const updateCategory = async (id, categoryData) => {
    try {
      await editCategory(id, categoryData);
      setCategories(prevCategories =>
        prevCategories.map(category =>
          category.id === id ? { ...category, ...categoryData } : category
        )
      );
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  const deleteCategory = async id => {
    try {
      await removeCategory(id);
      setCategories(prevCategories =>
        prevCategories.filter(category => category.id !== id)
      );
      return { success: true, id };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  return (
    <CategoriesContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={{
        categories,
        getCategories,
        error,
        loading,
        createCategory,
        updateCategory,
        deleteCategory
      }}
    >
      {children}
    </CategoriesContext.Provider>
  );
};

export const useCategories = () => {
  return useContext(CategoriesContext);
};
