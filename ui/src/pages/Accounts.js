import React, { useState, useEffect } from 'react';
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
import { useAccounts } from '../context';
import { AccountRow, Loading, Error } from '../components';

const Accounts = () => {
  const { accounts, deleteAccount, updateAccount, loading, error } =
    useAccounts();
  const [localAccounts, setLocalAccounts] = useState([]);
  const toast = useToast();

  useEffect(() => {
    setLocalAccounts(accounts);
  }, [accounts]);

  const handleDelete = async id => {
    await deleteAccount(id);
    toast({
      title: 'Account deleted.',
      description: 'Account has been successfully deleted.',
      status: 'info',
      duration: 1500,
      isClosable: true
    });
  };

  const handleSaveChanges = async (id, updatedAccount) => {
    try {
      await updateAccount(id, updatedAccount);
    } catch (err) {
      toast({
        title: 'Error updating account.',
        description: err.message,
        status: 'error',
        duration: 1500,
        isClosable: true
      });
    }
  };

  if (loading) return <Loading />;
  if (error) return <Error message={error.message} />;

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
            <Th padding="0.6rem">Account Name</Th>
            <Th padding="0.6rem">Account Category</Th>
            <Th padding="0.6rem">Off Budget</Th>
            <Th padding="0.6rem">Actions</Th>
          </Tr>
        </Thead>
        <Tbody>
          {localAccounts.length > 0 ? (
            localAccounts.map(account => (
              <AccountRow
                key={account.id}
                account={account}
                onDelete={handleDelete}
                onSave={handleSaveChanges}
              />
            ))
          ) : (
            <Tr>
              <Td
                padding="0.6rem"
                colSpan={4}
                textAlign="center"
                color="gray.500"
              >
                No accounts available
              </Td>
            </Tr>
          )}
        </Tbody>
      </Table>
    </Box>
  );
};

export default Accounts;
