import React, { useState, useEffect } from 'react';
import { Box, useTheme } from '@chakra-ui/react';
import { useBugdets, useGroups, useCategories } from '../context';
import {
  AddGroupModal,
  EditGroupModal,
  AddCategoryModal,
  EditCategoryModal,
  SetBudgetModal,
  YearMonthSelector,
  BudgetsTable
} from '../components/budgets';
import { Loading, Error } from '../components';

const Budgets = () => {
  const theme = useTheme();
  const primaryColor = theme.colors.primary;

  const { budgets, loading, error, getBudgets, createBudget } = useBugdets();
  const { createGroup, updateGroup, deleteGroup } = useGroups();
  const { createCategory, updateCategory, deleteCategory } = useCategories();

  // Months and current year setup
  const months = Array.from({ length: 12 }, (_, i) =>
    new Intl.DateTimeFormat('en-US', { month: 'long' }).format(new Date(0, i))
  );
  const currentYear = new Date().getFullYear();

  // State variables for managing groups, categories, and UI controls
  const [openGroups, setOpenGroups] = useState(
    budgets.reduce((acc, item) => ({ ...acc, [item.groupId]: true }), {})
  );
  const [selectedYear, setSelectedYear] = useState(currentYear);
  const [selectedMonth, setSelectedMonth] = useState(new Date().getMonth() + 1);
  const [hoveredRow, setHoveredRow] = useState(null);
  const [hoveredCategory, setHoveredCategory] = useState('');

  // Group management states
  const [isAddGroupOpen, setIsAddGroupOpen] = useState(false);
  const [isEditGroupOpen, setIsEditGroupOpen] = useState(false);
  const [formData, setFormData] = useState({ name: '', notes: '' });
  const [editFormData, setEditFormData] = useState({ name: '' });
  const [groupNameError, setGroupNameError] = useState('');
  const [selectedGroupId, setSelectedGroupId] = useState(null);

  // Category management states
  const [isAddCategoryOpen, setIsAddCategoryOpen] = useState(false);
  const [isEditCategoryOpen, setIsEditCategoryOpen] = useState(false);
  const [categoryFormData, setCategoryFormData] = useState({
    name: '',
    notes: ''
  });
  const [editCategoryFormData, setEditCategoryFormData] = useState({
    name: '',
    notes: ''
  });
  const [categoryNameError, setCategoryNameError] = useState('');
  const [selectedCategoryId, setSelectedCategoryId] = useState(null);

  // Budget management states
  const [isSetBudgetOpen, setIsSetBudgetOpen] = useState(false);
  const [budgetFormData, setBudgetFormData] = useState({ budgeted: 0 });

  // Fetch budgets when selectedYear or selectedMonth changes
  useEffect(() => {
    getBudgets(selectedYear, selectedMonth);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedYear, selectedMonth]);

  // Set initial open groups based on fetched budgets
  useEffect(() => {
    if (budgets.length) {
      const initialOpenGroups = budgets.reduce(
        (acc, item) => ({ ...acc, [item.groupId]: true }),
        {}
      );
      setOpenGroups(initialOpenGroups);
    }
  }, [budgets]);

  // Toggle group open/closed state
  const toggleGroup = groupId => {
    setOpenGroups(prev => ({ ...prev, [groupId]: !prev[groupId] }));
  };

  // Group and categorize budgets
  const groupedBudgets = budgets.reduce((acc, item) => {
    const { groupId, groupName, groupNotes, budgeted, spent } = item;
    if (!acc[groupId]) {
      acc[groupId] = {
        groupName,
        groupNotes,
        groupId,
        categories: [],
        budgeted: 0,
        spent: 0
      };
    }
    acc[groupId].categories.push(item);
    acc[groupId].budgeted += budgeted;
    acc[groupId].spent += spent;
    return acc;
  }, {});

  const groupedBudgetsArray = Object.values(groupedBudgets);

  // Handlers for form changes and submissions
  const handleChange = e => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    if (name === 'name' && value.trim()) setGroupNameError('');
  };

  const handleEditChange = e => {
    const { name, value } = e.target;
    setEditFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleCategoryChange = e => {
    const { name, value } = e.target;
    setCategoryFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleEditCategoryChange = e => {
    const { name, value } = e.target;
    setEditCategoryFormData(prev => ({ ...prev, [name]: value }));
  };

  // Group operations
  const handleAddGroupSubmit = async e => {
    e.preventDefault();
    if (!formData.name.trim()) {
      setGroupNameError('Group name is required');
      return;
    }

    const result = await createGroup({
      name: formData.name,
      notes: formData.notes || null
    });
    if (result.success) {
      setIsAddGroupOpen(false);
      setFormData({ name: '', notes: '' });
      await getBudgets(selectedYear, selectedMonth);
    } else {
      // eslint-disable-next-line no-console
      console.error('Error creating group:', result.error);
    }
  };

  const handleEditGroupSubmit = async e => {
    e.preventDefault();
    if (!editFormData.name.trim()) {
      setGroupNameError('Group name is required');
      return;
    }

    const result = await updateGroup(selectedGroupId, {
      name: editFormData.name,
      notes: editFormData.notes
    });
    if (result.success) {
      setIsEditGroupOpen(false);
      setEditFormData({ name: '', notes: '' });
      await getBudgets(selectedYear, selectedMonth);
    } else {
      // eslint-disable-next-line no-console
      console.error('Error updating group:', result.error);
    }
  };

  const handleDeleteGroup = async () => {
    const result = await deleteGroup(selectedGroupId);
    if (result.success) {
      setIsEditGroupOpen(false);
      setEditFormData({ name: '' });
      await getBudgets(selectedYear, selectedMonth);
    } else {
      // eslint-disable-next-line no-console
      console.error('Error deleting group:', result.error);
    }
  };

  // Category operations
  const handleAddCategorySubmit = async e => {
    e.preventDefault();
    if (!categoryFormData.name.trim()) {
      setCategoryNameError('Category name is required');
      return;
    }

    const result = await createCategory({
      name: categoryFormData.name,
      notes: categoryFormData.notes,
      groupId: selectedGroupId
    });

    if (result.success) {
      setIsAddCategoryOpen(false);
      setCategoryFormData({ name: '', notes: '' });
      await getBudgets(selectedYear, selectedMonth);
    } else {
      // eslint-disable-next-line no-console
      console.error('Error creating category:', result.error);
    }
  };

  const handleEditCategorySubmit = async e => {
    e.preventDefault();
    if (!editCategoryFormData.name.trim()) {
      setCategoryNameError('Category name is required');
      return;
    }

    const result = await updateCategory(selectedCategoryId, {
      groupId: selectedGroupId,
      name: editCategoryFormData.name,
      notes: editCategoryFormData.notes || null
    });

    if (result.success) {
      setIsEditCategoryOpen(false);
      setEditCategoryFormData({ name: '', notes: '' });
      await getBudgets(selectedYear, selectedMonth);
    } else {
      // eslint-disable-next-line no-console
      console.error('Error updating category:', result.error);
    }
  };

  const handleDeleteCategory = async () => {
    const result = await deleteCategory(selectedCategoryId);
    if (result.success) {
      setIsEditCategoryOpen(false);
      setEditCategoryFormData({ name: '', notes: '', groupId: '' });
      setSelectedCategoryId('');
      setSelectedGroupId('');
      await getBudgets(selectedYear, selectedMonth);
    } else {
      // eslint-disable-next-line no-console
      console.error('Error deleting category:', result.error);
    }
  };

  // Open edit modal with selected category data
  const openEditCategoryModal = (
    groupId,
    categoryId,
    categoryName,
    categoryNotes
  ) => {
    if (!categoryId) {
      setIsAddCategoryOpen(true);
      setSelectedGroupId(groupId);
      setCategoryFormData({ name: '', notes: '' });
    } else {
      setEditCategoryFormData({
        groupId,
        name: categoryName || '',
        notes: categoryNotes || ''
      });
      setSelectedCategoryId(categoryId);
      setSelectedGroupId(groupId);
      setIsEditCategoryOpen(true);
    }
  };

  // Budget operations
  const handleSetBudget = async e => {
    e.preventDefault();
    const result = await createBudget({
      month: selectedMonth,
      year: selectedYear,
      categoryId: selectedCategoryId,
      budgeted: parseFloat(budgetFormData.budgeted)
    });

    if (result.success) {
      setIsSetBudgetOpen(false);
      setBudgetFormData({ budgeted: 0 });
      setSelectedCategoryId('');
      await getBudgets(selectedYear, selectedMonth);
    } else {
      // eslint-disable-next-line no-console
      console.error('Error setting budget:', result.error);
    }
  };

  if (loading) return <Loading />;
  if (error) return <Error message={error.message} />;

  return (
    <Box>
      {/* Year and Month Selection */}
      <YearMonthSelector
        selectedYear={selectedYear}
        selectedMonth={selectedMonth}
        months={months}
        currentYear={currentYear}
        primaryColor={primaryColor}
        setSelectedYear={setSelectedYear}
        setSelectedMonth={setSelectedMonth}
        setIsAddGroupOpen={setIsAddGroupOpen}
      />

      {/* Add Group Modal */}
      <AddGroupModal
        isOpen={isAddGroupOpen}
        onClose={() => setIsAddGroupOpen(false)}
        formData={formData}
        groupNameError={groupNameError}
        handleChange={handleChange}
        handleAddGroupSubmit={handleAddGroupSubmit}
        primaryColor={primaryColor}
      />

      {/* Edit Group Modal */}
      <EditGroupModal
        isOpen={isEditGroupOpen}
        onClose={() => setIsEditGroupOpen(false)}
        editFormData={editFormData}
        groupNameError={groupNameError}
        handleEditChange={handleEditChange}
        handleEditGroupSubmit={handleEditGroupSubmit}
        handleDeleteGroup={handleDeleteGroup}
        primaryColor={primaryColor}
      />

      {/* Add Category Modal */}
      <AddCategoryModal
        isOpen={isAddCategoryOpen}
        onClose={() => setIsAddCategoryOpen(false)}
        categoryFormData={categoryFormData}
        categoryNameError={categoryNameError}
        handleCategoryChange={handleCategoryChange}
        handleAddCategorySubmit={handleAddCategorySubmit}
        primaryColor={primaryColor}
      />

      {/* Edit Category Modal */}
      <EditCategoryModal
        isOpen={isEditCategoryOpen}
        onClose={() => setIsEditCategoryOpen(false)}
        editCategoryFormData={editCategoryFormData}
        categoryNameError={categoryNameError}
        handleEditCategoryChange={handleEditCategoryChange}
        handleEditCategorySubmit={handleEditCategorySubmit}
        handleDeleteCategory={handleDeleteCategory}
        primaryColor={primaryColor}
      />

      {/* Set Budget Modal */}
      <SetBudgetModal
        isOpen={isSetBudgetOpen}
        onClose={() => setIsSetBudgetOpen(false)}
        budgetFormData={budgetFormData}
        setBudgetFormData={setBudgetFormData}
        handleSetBudget={handleSetBudget}
        primaryColor={primaryColor}
      />

      {/* Budgets Table */}
      {!loading && !error && (
        <BudgetsTable
          groupedBudgetsArray={groupedBudgetsArray}
          openGroups={openGroups}
          toggleGroup={toggleGroup}
          setHoveredRow={setHoveredRow}
          hoveredRow={hoveredRow}
          setEditFormData={setEditFormData}
          setSelectedGroupId={setSelectedGroupId}
          setIsEditGroupOpen={setIsEditGroupOpen}
          setIsAddCategoryOpen={setIsAddCategoryOpen}
          setCategoryFormData={setCategoryFormData}
          setHoveredCategory={setHoveredCategory}
          hoveredCategory={hoveredCategory}
          openEditCategoryModal={openEditCategoryModal}
          setIsSetBudgetOpen={setIsSetBudgetOpen}
          setSelectedCategoryId={setSelectedCategoryId}
          setBudgetFormData={setBudgetFormData}
          primaryColor={primaryColor}
        />
      )}
    </Box>
  );
};

export default Budgets;
