import React from 'react';
import { Box, Table, Thead, Tbody, Tr, Th, Td } from '@chakra-ui/react';
import TransactionRow from './TransactionRow';
import LoadingTransactions from '../LoadingTransactions';

const TransactionsTable = ({
  localTransactions,
  handleInputChange,
  handleSaveChanges,
  validationErrors,
  payees,
  categories,
  handleCheckboxChange,
  handleDelete,
  loading
}) => {
  if (loading) return <LoadingTransactions />;

  return (
    <Box overflowX="auto" borderRadius="md">
      <Table
        size="sm"
        variant="simple"
        bg="white"
        border="gray.300"
        borderRadius="md"
      >
        <Thead>
          <Tr>
            <Th padding="0.6rem">Name</Th>
            <Th padding="0.6rem">Payee</Th>
            <Th padding="0.6rem">Category</Th>
            <Th padding="0.6rem">Credit</Th>
            <Th padding="0.6rem">Debit</Th>
            <Th padding="0.6rem">Notes</Th>
            <Th padding="0.6rem">Cleared</Th>
            <Th padding="0.6rem">Actions</Th>
          </Tr>
        </Thead>
        <Tbody>
          {!loading && localTransactions.length > 0 ? (
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
                colSpan={8}
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
};

export default TransactionsTable;
