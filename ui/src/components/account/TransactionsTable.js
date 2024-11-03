import React from 'react';
import { Box, Table, Thead, Tbody, Tr, Th, Td } from '@chakra-ui/react';
import TransactionRow from './TransactionRow';

const TransactionsTable = ({
  localTransactions,
  handleInputChange,
  handleSaveChanges,
  validationErrors,
  payees,
  categories,
  handleCheckboxChange,
  handleDelete
}) => (
  <Box overflowX="auto" borderRadius="md">
    <Table variant="simple" bg="white" border="gray.300" borderRadius="md">
      <Thead>
        <Tr>
          <Th padding="0.5rem">Name</Th>
          <Th padding="0.5rem">Payee</Th>
          <Th padding="0.5rem">Category</Th>
          <Th padding="0.5rem">Credit</Th>
          <Th padding="0.5rem">Debit</Th>
          <Th padding="0.5rem">Notes</Th>
          <Th padding="0.5rem">Cleared At</Th>
          <Th padding="0.5rem">Actions</Th>
        </Tr>
      </Thead>
      <Tbody>
        {localTransactions.length > 0 ? (
          localTransactions.map((transaction, index) => (
            <TransactionRow
              key={transaction.id}
              transaction={transaction}
              index={index}
              handleInputChange={handleInputChange}
              handleSaveChanges={handleSaveChanges}
              validationErrors={validationErrors}
              payees={payees}
              categories={categories}
              handleCheckboxChange={handleCheckboxChange}
              handleDelete={handleDelete}
            />
          ))
        ) : (
          <Tr>
            <Td
              padding="0.6rem"
              colSpan={6}
              textAlign="center"
              color="gray.500"
            >
              No transactions available
            </Td>
          </Tr>
        )}
      </Tbody>
    </Table>
  </Box>
);

export default TransactionsTable;
