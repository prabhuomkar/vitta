/* eslint-disable no-nested-ternary */
import React from 'react';
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
import { useAccounts } from '../context';
import { AccountRow, Loading, Error } from '../components';

const Accounts = () => {
  const { accounts, deleteAccount, updateAccount, adapters, loading, error } =
    useAccounts();
  const toast = useToast();

  const handleDelete = async id => {
    await deleteAccount(id);
    toast({
      title: 'Account deleted.',
      description: 'Account has been successfully deleted.',
      status: 'success',
      duration: 1500,
      isClosable: true
    });
  };

  const handleSaveChanges = async (id, updatedAccount) => {
    const originalAccount = accounts.find(account => account.id === id);

    // function to check duplicate account name
    const duplicateAccount = accounts.some(
      account => account.name === updatedAccount.name && account.id !== id
    );

    if (duplicateAccount) {
      toast({
        title: 'Duplicate Account',
        description: `An account with the name "${updatedAccount.name}" already exists.`,
        status: 'warning',
        duration: 1500,
        isClosable: true
      });
      return;
    }

    const hasChanges = Object.keys(updatedAccount).some(
      key => updatedAccount[key] !== originalAccount[key]
    );

    if (!hasChanges) return;

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

  return (
    <Box overflowX="auto" borderRadius="md">
      {loading ? (
        <Loading />
      ) : error ? (
        <Error message={error.message} />
      ) : (
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
              <Th padding="0.6rem">Adapter</Th>
              <Th padding="0.6rem">Off Budget</Th>
              <Th padding="0.6rem">Actions</Th>
            </Tr>
          </Thead>
          <Tbody>
            {accounts.length > 0 ? (
              accounts.map(account => (
                <AccountRow
                  key={account.id}
                  account={account}
                  adapters={adapters}
                  onDelete={handleDelete}
                  onSave={handleSaveChanges}
                />
              ))
            ) : (
              <Tr>
                <Td
                  padding="0.6rem"
                  colSpan={5}
                  textAlign="center"
                  color="gray.500"
                >
                  <Box
                    display="flex"
                    flexDirection="column"
                    alignItems="center"
                  >
                    <Image
                      src={`${process.env.PUBLIC_URL}/assets/account.svg`}
                      alt="No accounts available"
                      width="200px"
                      height="200px"
                    />
                    <Text>No accounts available</Text>
                  </Box>
                </Td>
              </Tr>
            )}
          </Tbody>
        </Table>
      )}
    </Box>
  );
};

export default Accounts;
