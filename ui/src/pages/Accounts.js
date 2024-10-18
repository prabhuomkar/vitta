import React, { useState, useEffect, useContext, useCallback } from 'react';
import {
  Box,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Input,
  Select,
  Checkbox,
  Button,
  useToast,
  Skeleton,
  Stack
} from '@chakra-ui/react';
import { debounce } from 'lodash';
import { AccountsContext } from '../context';
import { ACCOUNT_CATEGORIES } from '../components/common/constants';

const Accounts = () => {
  const { accounts, loading, error, deleteAccount, updateAccount } =
    useContext(AccountsContext);
  const [localAccounts, setLocalAccounts] = useState([]);
  const [errors, setErrors] = useState({});
  const toast = useToast();

  useEffect(() => {
    setLocalAccounts(accounts);
  }, [accounts]);

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const debouncedUpdateAccount = useCallback(
    debounce((id, updatedAccount) => {
      updateAccount(id, updatedAccount).catch(err => {
        toast({
          title: 'Error updating account.',
          description: err.message,
          status: 'error',
          duration: 3000,
          isClosable: true
        });
      });
    }, 500),
    []
  );

  const handleDelete = async id => {
    await deleteAccount(id);
    toast({
      title: 'Account deleted.',
      description: `Account has been successfully deleted.`,
      status: 'info',
      duration: 3000,
      isClosable: true
    });
  };

  const handleFieldChange = (id, field, value) => {
    const updatedAccounts = localAccounts.map(account =>
      account.id === id ? { ...account, [field]: value } : account
    );
    setLocalAccounts(updatedAccounts);

    // Update the errors state if the field is empty
    setErrors(prevErrors => ({
      ...prevErrors,
      [id]: {
        ...prevErrors[id],
        [field]: !value // Set to true if the field is empty, otherwise false
      }
    }));

    // Validate required fields before making the API call
    const updatedAccount = updatedAccounts.find(account => account.id === id);
    const hasErrors = !updatedAccount.name || !updatedAccount.category; // Check if required fields are empty

    if (!hasErrors) {
      // Only call the debounced function if there are no validation errors
      debouncedUpdateAccount(id, updatedAccount);
    }
  };

  const isFieldInvalid = (id, field) => errors[id]?.[field];

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

  if (error) return <Box>Error loading accounts: {error.message}</Box>;

  return (
    <Box overflowX="auto" borderRadius="md">
      <Table variant="simple" bg="white" border="gray.300" borderRadius="md">
        <Thead>
          <Tr>
            <Th>Account Name</Th>
            <Th>Account Category</Th>
            <Th>Off Budget</Th>
            <Th>Actions</Th>
          </Tr>
        </Thead>
        <Tbody>
          {localAccounts.length > 0 ? (
            localAccounts.map(account => (
              <Tr key={`${account.id}`}>
                <Td>
                  <Input
                    value={account.name || ''}
                    onChange={e =>
                      handleFieldChange(account.id, 'name', e.target.value)
                    }
                    placeholder="Account Name"
                    isInvalid={isFieldInvalid(account.id, 'name')}
                    errorBorderColor="red.300"
                  />
                </Td>
                <Td>
                  <Select
                    value={account.category || ''}
                    onChange={e =>
                      handleFieldChange(account.id, 'category', e.target.value)
                    }
                    placeholder="Select Category"
                    isInvalid={isFieldInvalid(account.id, 'category')}
                    errorBorderColor="red.300"
                  >
                    {ACCOUNT_CATEGORIES.map(cat => (
                      <option key={cat} value={cat}>
                        {cat}
                      </option>
                    ))}
                  </Select>
                </Td>
                <Td>
                  <Checkbox
                    isChecked={account.offBudget || false}
                    onChange={e =>
                      handleFieldChange(
                        account.id,
                        'offBudget',
                        e.target.checked
                      )
                    }
                  >
                    Off Budget
                  </Checkbox>
                </Td>
                <Td>
                  <Button
                    colorScheme="red"
                    onClick={() => handleDelete(account.id)}
                  >
                    Delete
                  </Button>
                </Td>
              </Tr>
            ))
          ) : (
            <Tr>
              <Td colSpan={4} textAlign="center" color="gray.500">
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
