import React from 'react';
import {
  Box,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Collapse,
  Flex,
  IconButton,
  Tooltip,
  Text,
  Image,
  Button,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  Portal
} from '@chakra-ui/react';
import { ChevronDownIcon, ChevronUpIcon } from '@chakra-ui/icons';
import { GoKebabHorizontal } from 'react-icons/go';
import LoadingBudgets from '../LoadingBudgets';
import { formatCurrency } from '../../utils/formatCurrency';

const BudgetsTable = ({
  groupedBudgetsArray,
  openGroups,
  toggleGroup,
  setHoveredRow,
  hoveredRow,
  setEditFormData,
  setSelectedGroupId,
  setIsEditGroupOpen,
  setIsAddCategoryOpen,
  setCategoryFormData,
  setHoveredCategory,
  hoveredCategory,
  openEditCategoryModal,
  setIsSetBudgetOpen,
  setSelectedCategoryId,
  setBudgetFormData,
  handleDeleteGroup,
  handleDeleteCategory,
  loading
}) => {
  if (loading) return <LoadingBudgets />;

  return (
    <Box overflowX="auto">
      <Table
        size="sm"
        variant="simple"
        bg="white"
        border="gray.300"
        borderRadius="md"
      >
        <Thead>
          <Tr>
            <Th padding="0.6rem">Group / Category</Th>
            <Th padding="0.6rem">Budgeted</Th>
            <Th padding="0.6rem">Spent</Th>
            <Th padding="0.6rem">Balance</Th>
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
                <Box display="flex" flexDirection="column" alignItems="center">
                  <Image
                    src={`${process.env.PUBLIC_URL}/assets/savings.svg`}
                    alt="No budgets available"
                    width="200px"
                    height="200px"
                  />
                  <Text>No budgets available</Text>
                </Box>
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
                  <Td width="40%" fontWeight="medium">
                    <Flex alignItems="center">
                      <IconButton
                        size="xs"
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
                        <Text cursor="pointer">{group.groupName}</Text>
                      </Tooltip>
                      {hoveredRow === group.groupId && (
                        <Menu>
                          <MenuButton
                            as={IconButton}
                            size="xs"
                            ml={2}
                            aria-label="Options"
                            icon={<GoKebabHorizontal fontSize="14px" />}
                            variant="outline"
                          />
                          <Portal>
                            <MenuList fontSize="sm">
                              <MenuItem
                                onClick={() => {
                                  setIsAddCategoryOpen(true);
                                  setSelectedGroupId(group.groupId);
                                  setCategoryFormData({ name: '', notes: '' });
                                }}
                              >
                                Add Category
                              </MenuItem>
                              <MenuItem
                                onClick={() => {
                                  setEditFormData({
                                    name: group.groupName,
                                    notes: group.groupNotes
                                  });
                                  setSelectedGroupId(group.groupId);
                                  setIsEditGroupOpen(true);
                                }}
                              >
                                Edit Group
                              </MenuItem>
                              <MenuItem
                                onClick={() => handleDeleteGroup(group.groupId)}
                              >
                                Delete Group
                              </MenuItem>
                            </MenuList>
                          </Portal>
                        </Menu>
                      )}
                    </Flex>
                  </Td>
                  <Td width="20%" fontWeight="medium">
                    {formatCurrency(group.budgeted)}
                  </Td>
                  <Td width="20%" fontWeight="medium">
                    {formatCurrency(group.spent)}
                  </Td>
                  <Td width="20%" fontWeight="medium">
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
                                  cursor="pointer"
                                  onMouseEnter={() =>
                                    setHoveredCategory(category.categoryId)
                                  }
                                  onMouseLeave={() => setHoveredCategory('')}
                                >
                                  <Flex marginLeft="32px">
                                    <Tooltip
                                      hasArrow
                                      openDelay={500}
                                      label={category.categoryNotes || ''}
                                    >
                                      <Text cursor="pointer" mt={1}>
                                        {category.categoryName || (
                                          <Button
                                            size="xs"
                                            onClick={() =>
                                              openEditCategoryModal(
                                                group.groupId,
                                                category.categoryId,
                                                category.categoryName,
                                                category.categoryNotes
                                              )
                                            }
                                          >
                                            Add Category
                                          </Button>
                                        )}
                                      </Text>
                                    </Tooltip>
                                    {category.categoryId && (
                                      <Menu>
                                        <MenuButton
                                          as={IconButton}
                                          size="xs"
                                          ml={2}
                                          aria-label="Options"
                                          icon={
                                            <GoKebabHorizontal fontSize="14px" />
                                          }
                                          variant="outline"
                                          visibility={
                                            hoveredCategory ===
                                            category.categoryId
                                              ? 'visible'
                                              : 'hidden'
                                          }
                                        />
                                        <Portal>
                                          <MenuList fontSize="sm">
                                            <MenuItem
                                              onClick={() => {
                                                setIsSetBudgetOpen(true);
                                                setSelectedCategoryId(
                                                  category.categoryId
                                                );
                                                setBudgetFormData({
                                                  budgeted: category.budgeted
                                                });
                                              }}
                                            >
                                              Set Budget
                                            </MenuItem>
                                            <MenuItem
                                              onClick={() =>
                                                openEditCategoryModal(
                                                  group.groupId,
                                                  category.categoryId,
                                                  category.categoryName,
                                                  category.categoryNotes
                                                )
                                              }
                                            >
                                              Edit Category
                                            </MenuItem>
                                            <MenuItem
                                              onClick={() =>
                                                handleDeleteCategory(
                                                  category.categoryId
                                                )
                                              }
                                            >
                                              Delete Category
                                            </MenuItem>
                                          </MenuList>
                                        </Portal>
                                      </Menu>
                                    )}
                                  </Flex>
                                </Td>
                                <Td width="20%">
                                  {formatCurrency(category.budgeted)}
                                </Td>
                                <Td width="20%">
                                  {formatCurrency(category.spent)}
                                </Td>
                                <Td width="20%">
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
  );
};

export default React.memo(BudgetsTable);
