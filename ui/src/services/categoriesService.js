import {
  getCategories,
  createCategory,
  updateCatrgory,
  deleteCategory
} from '../api/categoriesApi';

export const fetchCategories = async () => {
  try {
    const response = await getCategories();
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error fetching categories:', error);
    throw error;
  }
};

export const addCategory = async categoryData => {
  try {
    const response = await createCategory(categoryData);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error creating category:', error);
    throw error;
  }
};

export const editCategory = async (id, categoryData) => {
  try {
    const response = await updateCatrgory(id, categoryData);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error updating category:', error);
    throw error;
  }
};

export const removeCategory = async id => {
  try {
    const response = await deleteCategory(id);
    return response.data;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.error('Error deleting category:', error);
    throw error;
  }
};
