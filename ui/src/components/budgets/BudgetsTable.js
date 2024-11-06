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
  Button
} from '@chakra-ui/react';
import {
  ChevronDownIcon,
  ChevronUpIcon,
  HamburgerIcon
} from '@chakra-ui/icons';
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
  primaryColor
}) => {
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
                          size="xs"
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

export default BudgetsTable;
