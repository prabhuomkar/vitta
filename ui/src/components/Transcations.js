import React, { useEffect } from 'react';
import {
  Box,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  useToast,
  Skeleton,
  Stack
} from '@chakra-ui/react';
import { useParams } from 'react-router-dom';
import { useTransactions } from '../context';

const Transactions = () => {
  const { accountId } = useParams();
  const { transactions, loading, error, loadTransactions } = useTransactions();
  const toast = useToast();

  // Load transactions when the component mounts or accountId changes
  useEffect(() => {
    const fetchTransactions = async () => {
      if (accountId) {
        await loadTransactions(accountId); // Fetch transactions using accountId
      }
    };

    fetchTransactions();
  }, [accountId, loadTransactions]);

  if (loading) {
    return (
      <Box bg="white" padding="1rem" borderRadius="md">
        <Stack>
          <Skeleton height="60px" />
          <Skeleton height="60px" />
        </Stack>
      </Box>
    );
  }

  if (error) {
    toast({
      title: 'Error loading transactions.',
      description: error.message,
      status: 'error',
      duration: 3000,
      isClosable: true
    });
    return <Box>Error loading transactions: {error.message}</Box>;
  }

  return (
    <Box overflowX="auto" borderRadius="md">
      <Table variant="simple" bg="white" border="gray.300" borderRadius="md">
        <Thead>
          <Tr>
            <Th padding="0.6rem">Name</Th>
            <Th padding="0.6rem">Credited</Th>
            <Th padding="0.6rem">Debit</Th>
            <Th padding="0.6rem">Notes</Th>
          </Tr>
        </Thead>
        <Tbody>
          {transactions.length > 0 ? (
            transactions.map(transaction => (
              <Tr key={`${transaction.id}`}>
                <Td padding="0.6rem">{transaction.name || 'N/A'}</Td>
                <Td padding="0.6rem">{transaction.credit || 'N/A'}</Td>
                <Td padding="0.6rem">{transaction.debit || 'N/A'}</Td>
                <Td padding="0.6rem">{transaction.notes || 'N/A'}</Td>
              </Tr>
            ))
          ) : (
            <Tr>
              <Td
                padding="0.6rem"
                colSpan={4}
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

export default Transactions;
