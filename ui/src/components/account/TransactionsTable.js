import React from 'react';
import {
  Box,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  useToast
} from '@chakra-ui/react';
import LoadingTransactions from '../LoadingTransactions';
import TransactionRow from './TransactionRow';
import { useTransactions, usePayees, useCategories } from '../../context';

const TransactionsTable = () => {
  const toast = useToast();
  const { transactions, loading, updateTransaction, deleteTransaction } =
    useTransactions();
  const { payees } = usePayees();
  const { categories } = useCategories();

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
          {transactions.length > 0 ? (
            transactions.map(transaction => (
              <TransactionRow
                key={transaction.id}
                transaction={transaction}
                payees={payees}
                categories={categories}
                updateTransaction={updateTransaction}
                deleteTransaction={deleteTransaction}
                toast={toast}
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
