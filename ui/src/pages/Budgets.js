import React, { useState, useEffect } from 'react';
import {
  Box,
  Flex,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Collapse,
  IconButton,
  Select,
  Text,
  Heading,
  Button,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  FormControl,
  FormLabel,
  Input,
  FormErrorMessage,
  Tooltip,
  Textarea,
  useTheme
} from '@chakra-ui/react';
import {
  ChevronDownIcon,
  ChevronUpIcon,
  HamburgerIcon
} from '@chakra-ui/icons';
import { useBugdets, useGroups, useCategories } from '../context';
import { formatCurrency } from '../utils/formatCurrency';
import { Loading, Error } from '../components';

const Budgets = () => {
  const theme = useTheme();
  const primaryColor = theme.colors.primary;

  const { budgets, loading, error, getBudgets, createBudget } = useBugdets();
  const { createGroup, updateGroup, deleteGroup } = useGroups();
  const { createCategory, updateCategory, deleteCategory } = useCategories();
  const months = Array.from({ length: 12 }, (_, i) =>
    new Intl.DateTimeFormat('en-US', { month: 'long' }).format(new Date(0, i))
  );
  const currentYear = new Date().getFullYear();
  const [openGroups, setOpenGroups] = useState(
    budgets.reduce((acc, item) => {
      acc[item.groupId] = true;
      return acc;
    }, {})
  );

  const [selectedYear, setSelectedYear] = useState(currentYear);
  const [selectedMonth, setSelectedMonth] = useState(new Date().getMonth() + 1);
  const [isAddGroupOpen, setIsAddGroupOpen] = useState(false);
  const [isEditGroupOpen, setIsEditGroupOpen] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    notes: ''
  });

  const [isEditCategoryOpen, setIsEditCategoryOpen] = useState(false);
  const [selectedCategoryId, setSelectedCategoryId] = useState(null);
  const [editCategoryFormData, setEditCategoryFormData] = useState({
    name: '',
    notes: ''
  });

  const [isAddCategoryOpen, setIsAddCategoryOpen] = useState(false);
  const [categoryFormData, setCategoryFormData] = useState({
    name: '',
    notes: ''
  });
  const [categoryNameError, setCategoryNameError] = useState('');

  const [editFormData, setEditFormData] = useState({
    name: ''
  });
  const [groupNameError, setGroupNameError] = useState('');
  const [hoveredRow, setHoveredRow] = useState(null);
  const [selectedGroupId, setSelectedGroupId] = useState(null);

  const [isSetBudgetOpen, setIsSetBudgetOpen] = useState(false);
  const [budgetFormData, setBudgetFormData] = useState({
    budgeted: 0
  });
  const [hoveredCategory, setHoveredCategory] = useState('');

  useEffect(() => {
    const fetchBudgets = async () => {
      await getBudgets(selectedYear, selectedMonth);
    };

    fetchBudgets();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedYear, selectedMonth]);

  useEffect(() => {
    if (budgets.length > 0) {
      const initialOpenGroups = budgets.reduce((acc, item) => {
        acc[item.groupId] = true; // Set all groups to open by default
        return acc;
      }, {});
      setOpenGroups(initialOpenGroups);
    }
  }, [budgets]);

  const toggleGroup = groupId => {
    setOpenGroups(prev => ({
      ...prev,
      [groupId]: !prev[groupId]
    }));
  };

  const groupedBudgets = budgets.reduce((acc, item) => {
    const { groupId } = item;
    if (!acc[groupId]) {
      acc[groupId] = {
        groupName: item.groupName,
        groupNotes: item.groupNotes,
        groupId,
        categories: [],
        budgeted: 0,
        spent: 0
      };
    }
    acc[groupId].categories.push(item);
    acc[groupId].budgeted += item.budgeted;
    acc[groupId].spent += item.spent;
    return acc;
  }, {});

  const groupedBudgetsArray = Object.values(groupedBudgets);

  const handleChange = e => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    if (name === 'name' && value.trim() !== '') {
      setGroupNameError('');
    }
  };

  const handleEditChange = e => {
    const { name, value } = e.target;
    setEditFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

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

  // Function to handle category form change
  const handleCategoryChange = e => {
    const { name, value } = e.target;

    setCategoryFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // Function to handle adding a new category
  const handleAddCategorySubmit = async e => {
    e.preventDefault();
    if (!categoryFormData.name.trim()) {
      setCategoryNameError('Category name is required');
      return;
    }

    const result = await createCategory({
      name: categoryFormData.name,
      notes: categoryFormData.notes,
      groupId: selectedGroupId // Use the selected group ID
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

  // Function to handle category form change in edit modal
  const handleEditCategoryChange = e => {
    const { name, value } = e.target;
    setEditCategoryFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // Function to submit the edited category
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
      setIsEditCategoryOpen(false); // Close the edit modal
      setEditCategoryFormData({ name: '', notes: '', groupId: '' });
      setSelectedCategoryId(''); // Clear selected category
      setSelectedGroupId(''); // Clear selected group
      await getBudgets(selectedYear, selectedMonth);
    } else {
      // eslint-disable-next-line no-console
      console.error('Error deleting category:', result.error);
    }
  };

  // Open the edit category modal and set selected category data
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
      <Box mb={4}>
        <Flex
          justifyContent="space-between"
          alignItems="center"
          flexDirection={{ base: 'column', md: 'row' }}
        >
          <Heading mb={{ base: 2, md: 0 }}>{selectedYear}</Heading>
          <Flex
            alignItems="center"
            flexDirection={{ base: 'column', md: 'row' }}
            width={{ base: '100%', md: 'auto' }}
          >
            <Select
              value={selectedMonth}
              onChange={e => setSelectedMonth(Number(e.target.value))}
              placeholder="Select Month"
              width={{ base: '100%', md: 'auto' }}
              backgroundColor="white"
              m={2}
            >
              {months.map((month, index) => (
                // eslint-disable-next-line react/no-array-index-key
                <option key={index} value={index + 1}>
                  {month}
                </option>
              ))}
            </Select>
            <Select
              value={selectedYear}
              onChange={e => setSelectedYear(Number(e.target.value))}
              placeholder="Select Year"
              width={{ base: '100%', md: 'auto' }}
              backgroundColor="white"
              m={2}
            >
              {Array.from({ length: 5 }, (_, i) => (
                <option key={i} value={currentYear - 2 + i}>
                  {currentYear - 2 + i}
                </option>
              ))}
            </Select>

            <Button
              colorScheme={primaryColor}
              bg={primaryColor}
              onClick={() => setIsAddGroupOpen(true)}
              width={{ base: '100%', md: 'auto' }}
              ml={{ base: 0, md: 2 }}
              mt={{ base: 2, md: 0 }}
            >
              Add Group
            </Button>
          </Flex>
        </Flex>
      </Box>

      {/* Modal for adding group */}
      <Modal isOpen={isAddGroupOpen} onClose={() => setIsAddGroupOpen(false)}>
        <ModalOverlay />
        <ModalContent mx={{ base: '4', md: '0' }}>
          <ModalHeader>Add Group</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <form onSubmit={handleAddGroupSubmit}>
              <FormControl mb={4} isInvalid={!!groupNameError}>
                <FormLabel>Group Name</FormLabel>
                <Input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  placeholder="Enter Group Name"
                />
                {groupNameError && (
                  <FormErrorMessage>{groupNameError}</FormErrorMessage>
                )}
              </FormControl>
              <FormControl mb={4}>
                <FormLabel>Notes</FormLabel>
                <Textarea
                  name="notes"
                  value={formData.notes}
                  onChange={handleChange}
                  placeholder="Enter Notes"
                />
              </FormControl>
              <ModalFooter margin="auto" mb={2} p={0}>
                <Button
                  colorScheme={primaryColor}
                  bg={primaryColor}
                  type="submit"
                >
                  Add Group
                </Button>
              </ModalFooter>
            </form>
          </ModalBody>
        </ModalContent>
      </Modal>

      {/* Modal for editing group */}
      <Modal isOpen={isEditGroupOpen} onClose={() => setIsEditGroupOpen(false)}>
        <ModalOverlay />
        <ModalContent mx={{ base: '4', md: '0' }}>
          <ModalHeader>Edit Group</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <form onSubmit={handleEditGroupSubmit}>
              <FormControl mb={4} isInvalid={!!groupNameError}>
                <FormLabel>Group Name</FormLabel>
                <Input
                  type="text"
                  name="name"
                  value={editFormData.name}
                  onChange={handleEditChange}
                  placeholder="Enter Group Name"
                />
                {groupNameError && (
                  <FormErrorMessage>{groupNameError}</FormErrorMessage>
                )}
              </FormControl>
              <FormControl mb={4}>
                <FormLabel>Notes</FormLabel>
                <Textarea
                  name="notes"
                  value={editFormData.notes}
                  onChange={handleEditChange}
                  placeholder="Enter Notes"
                />
              </FormControl>
              <ModalFooter margin="auto" mb={2} p={0}>
                <Button
                  colorScheme={primaryColor}
                  bg={primaryColor}
                  type="submit"
                  mr={3}
                >
                  Update
                </Button>
                <Button onClick={handleDeleteGroup}>Delete</Button>
              </ModalFooter>
            </form>
          </ModalBody>
        </ModalContent>
      </Modal>

      {/* Modal for adding category */}
      <Modal
        isOpen={isAddCategoryOpen}
        onClose={() => setIsAddCategoryOpen(false)}
      >
        <ModalOverlay />
        <ModalContent mx={{ base: '4', md: '0' }}>
          <ModalHeader>Add Category</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <form onSubmit={handleAddCategorySubmit}>
              <FormControl mb={4} isInvalid={!!categoryNameError}>
                <FormLabel>Category Name</FormLabel>
                <Input
                  type="text"
                  name="name"
                  value={categoryFormData.name}
                  onChange={handleCategoryChange}
                  placeholder="Enter Category Name"
                />
                {categoryNameError && (
                  <FormErrorMessage>{categoryNameError}</FormErrorMessage>
                )}
              </FormControl>
              <FormControl mb={4}>
                <FormLabel>Notes</FormLabel>
                <Textarea
                  type="text"
                  name="notes"
                  value={categoryFormData.notes}
                  onChange={handleCategoryChange}
                  placeholder="Enter Notes"
                />
              </FormControl>
              <ModalFooter margin="auto" mb={2} p={0}>
                <Button
                  colorScheme={primaryColor}
                  bg={primaryColor}
                  type="submit"
                >
                  Add Category
                </Button>
              </ModalFooter>
            </form>
          </ModalBody>
        </ModalContent>
      </Modal>

      {/* Edit Category Modal */}
      <Modal
        isOpen={isEditCategoryOpen}
        onClose={() => setIsEditCategoryOpen(false)}
      >
        <ModalOverlay />
        <ModalContent mx={{ base: '4', md: '0' }}>
          <ModalHeader>Edit Category</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <form onSubmit={handleEditCategorySubmit}>
              <FormControl mb={4} isInvalid={!!categoryNameError}>
                <FormLabel>Category Name</FormLabel>
                <Input
                  type="text"
                  name="name"
                  value={editCategoryFormData.name}
                  onChange={handleEditCategoryChange}
                  placeholder="Enter Category Name"
                />
                {categoryNameError && (
                  <FormErrorMessage>{categoryNameError}</FormErrorMessage>
                )}
              </FormControl>
              <FormControl mb={4}>
                <FormLabel>Notes</FormLabel>
                <Textarea
                  type="text"
                  name="notes"
                  value={editCategoryFormData.notes}
                  onChange={handleEditCategoryChange}
                  placeholder="Enter Notes"
                />
              </FormControl>
              <ModalFooter margin="auto" mb={2} p={0}>
                <Button
                  colorScheme={primaryColor}
                  bg={primaryColor}
                  type="submit"
                  mr={3}
                >
                  Update
                </Button>
                <Button onClick={handleDeleteCategory}>Delete</Button>
              </ModalFooter>
            </form>
          </ModalBody>
        </ModalContent>
      </Modal>

      <Modal isOpen={isSetBudgetOpen} onClose={() => setIsSetBudgetOpen(false)}>
        <ModalOverlay />
        <ModalContent mx={{ base: '4', md: '0' }}>
          <ModalHeader>Add Budget</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <form onSubmit={handleSetBudget}>
              <FormControl mb={4}>
                <FormLabel>Budget</FormLabel>
                <Input
                  type="number"
                  name="budgeted"
                  value={budgetFormData.budgeted}
                  onChange={e =>
                    setBudgetFormData({ budgeted: e.target.value })
                  }
                  placeholder="Enter Budget Amount"
                />
              </FormControl>
              <ModalFooter margin="auto" mb={2} p={0}>
                <Button
                  colorScheme={primaryColor}
                  bg={primaryColor}
                  type="submit"
                >
                  Set Budget
                </Button>
              </ModalFooter>
            </form>
          </ModalBody>
        </ModalContent>
      </Modal>

      {/* Table for budgets */}
      {!loading && !error && (
        <Box overflowX="auto">
          <Table
            variant="simple"
            bg="white"
            border="gray.300"
            borderRadius="md"
          >
            <Thead>
              <Tr>
                <Th>Group / Category</Th>
                <Th>Budgeted</Th>
                <Th>Spent</Th>
                <Th>Balance</Th>
              </Tr>
            </Thead>
            <Tbody>
              {groupedBudgetsArray.length === 0 ? (
                <Tr>
                  <Td
                    padding="0.6rem"
                    colSpan={4}
                    textAlign="center"
                    color="gray.500"
                  >
                    <Text>No budgets available</Text>
                  </Td>
                </Tr>
              ) : (
                groupedBudgetsArray.map(group => (
                  <React.Fragment key={group.groupId}>
                    {/* Group Row */}
                    <Tr
                      onMouseEnter={() => setHoveredRow(group.groupId)}
                      onMouseLeave={() => setHoveredRow(null)}
                    >
                      <Td width="40%">
                        <Flex alignItems="center">
                          <IconButton
                            size="sm"
                            variant="outline"
                            onClick={() => toggleGroup(group.groupId)}
                            icon={
                              openGroups[group.groupId] ? (
                                <ChevronUpIcon boxSize={5} />
                              ) : (
                                <ChevronDownIcon boxSize={5} />
                              )
                            }
                            aria-label={`Toggle ${group.groupName}`}
                            mr={2}
                          />
                          <Tooltip
                            hasArrow
                            openDelay={500}
                            label={group.groupNotes}
                          >
                            <Text
                              cursor="pointer"
                              onClick={() => {
                                setEditFormData({
                                  name: group.groupName,
                                  notes: group.groupNotes
                                });
                                setSelectedGroupId(group.groupId);
                                setIsEditGroupOpen(true);
                              }}
                            >
                              {group.groupName}
                            </Text>
                          </Tooltip>
                          {hoveredRow === group.groupId && (
                            <IconButton
                              size="sm"
                              variant="outline"
                              icon={<HamburgerIcon />}
                              onClick={() => {
                                setIsAddCategoryOpen(true);
                                setSelectedGroupId(group.groupId);
                                setCategoryFormData({ name: '', notes: '' });
                              }}
                              ml={2}
                            />
                          )}
                        </Flex>
                      </Td>
                      <Td width="20%">{formatCurrency(group.budgeted)}</Td>
                      <Td width="20%">{formatCurrency(group.spent)}</Td>
                      <Td width="20%">
                        {formatCurrency(group.budgeted + group.spent)}
                      </Td>
                    </Tr>
                    {/* Collapsible Categories */}
                    <Tr>
                      <Td colSpan={4} padding={0}>
                        <Collapse in={openGroups[group.groupId]} animateOpacity>
                          <Box>
                            <Table size="sm">
                              <Tbody>
                                {group.categories.map(category => (
                                  <Tr key={category.categoryId}>
                                    <Td
                                      width="40%"
                                      textAlign="center"
                                      cursor="pointer"
                                      onMouseEnter={() =>
                                        setHoveredCategory(category.categoryId)
                                      }
                                      onMouseLeave={() =>
                                        setHoveredCategory('')
                                      }
                                    >
                                      <Flex
                                        alignItems="center"
                                        justifyContent="center"
                                      >
                                        <Tooltip
                                          hasArrow
                                          openDelay={500}
                                          label={category.categoryNotes || ''}
                                        >
                                          <Text
                                            cursor="pointer"
                                            onClick={() =>
                                              openEditCategoryModal(
                                                group.groupId,
                                                category.categoryId,
                                                category.categoryName,
                                                category.categoryNotes
                                              )
                                            }
                                          >
                                            {category.categoryName || (
                                              <Button
                                                colorScheme={primaryColor}
                                                bg={primaryColor}
                                                size="xs"
                                              >
                                                Add Category
                                              </Button>
                                            )}
                                          </Text>
                                        </Tooltip>
                                        {category.categoryId && (
                                          <IconButton
                                            size="xs"
                                            variant="outline"
                                            icon={<HamburgerIcon />}
                                            onClick={() => {
                                              setIsSetBudgetOpen(true);
                                              setSelectedCategoryId(
                                                category.categoryId
                                              );
                                              setBudgetFormData({
                                                budgeted: category.budgeted
                                              });
                                            }}
                                            ml={2}
                                            visibility={
                                              hoveredCategory ===
                                              category.categoryId
                                                ? 'visible'
                                                : 'hidden'
                                            }
                                          />
                                        )}
                                      </Flex>
                                    </Td>
                                    <Td width="20%" textAlign="center">
                                      {formatCurrency(category.budgeted)}
                                    </Td>
                                    <Td width="20%" textAlign="center">
                                      {formatCurrency(category.spent)}
                                    </Td>
                                    <Td width="20%" textAlign="center">
                                      {formatCurrency(
                                        category.budgeted + category.spent
                                      )}
                                    </Td>
                                  </Tr>
                                ))}
                              </Tbody>
                            </Table>
                          </Box>
                        </Collapse>
                      </Td>
                    </Tr>
                  </React.Fragment>
                ))
              )}
            </Tbody>
          </Table>
        </Box>
      )}
    </Box>
  );
};

export default Budgets;
