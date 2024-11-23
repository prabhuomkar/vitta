import React, { useEffect } from 'react';
import {
  Box,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Image,
  Text,
  useToast
} from '@chakra-ui/react';
import LoadingTransactions from '../LoadingTransactions';
import TransactionRow from './TransactionRow';
import { useTransactions, usePayees, useCategories } from '../../context';

const TransactionsTable = ({ getAccountById }) => {
  const toast = useToast();
  const {
    transactions,
    page,
    totalPages,
    goToPreviousPage,
    loading,
    updateTransaction,
    deleteTransaction
  } = useTransactions();
  const { payees } = usePayees();
  const { categories } = useCategories();

  useEffect(() => {
    if (transactions.length === 0 && page > 1 && totalPages > 1) {
      goToPreviousPage();
    }
  }, [transactions.length, page, totalPages, goToPreviousPage]);

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
                getAccountById={getAccountById}
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
                <Box display="flex" flexDirection="column" alignItems="center">
                  <Image
                    src={`${process.env.PUBLIC_URL}/assets/online_payments.svg`}
                    alt="No transactions available"
                    width="200px"
                    height="200px"
                  />
                  <Text>No transactions available</Text>
                </Box>
              </Td>
            </Tr>
          )}
        </Tbody>
      </Table>
    </Box>
  );
};

export default TransactionsTable;
