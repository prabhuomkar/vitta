import React, { useState, useEffect, useContext } from 'react';
import {
  Box,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Button,
  useToast,
  Skeleton,
  Stack,
  useTheme,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Input,
  Select,
  Checkbox
} from '@chakra-ui/react';
import { AccountsContext } from '../context';
import { ACCOUNT_CATEGORIES } from '../components/common/constants';

const Accounts = () => {
  const theme = useTheme();
  const primaryColor = theme.colors.primary;
  const {
    accounts,
    getAccounts,
    loading,
    error,
    deleteAccount,
    updateAccount
  } = useContext(AccountsContext);
  const [localAccounts, setLocalAccounts] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedAccount, setSelectedAccount] = useState(null);
  const [nameError, setNameError] = useState('');
  const [categoryError, setCategoryError] = useState('');
  const toast = useToast();

  useEffect(() => {
    setLocalAccounts(accounts);
  }, [accounts]);

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

  const handleEdit = account => {
    setSelectedAccount(account);
    setIsModalOpen(true);
    setNameError('');
    setCategoryError('');
  };

  const handleSave = async () => {
    if (!selectedAccount) return;

    setNameError('');
    setCategoryError('');

    let isValid = true;

    if (!selectedAccount.name.trim()) {
      setNameError('Account name cannot be empty.');
      isValid = false;
    }

    if (!selectedAccount.category) {
      setCategoryError('Please select a category.');
      isValid = false;
    }

    if (!isValid) {
      return;
    }

    try {
      await updateAccount(selectedAccount.id, {
        name: selectedAccount.name,
        category: selectedAccount.category,
        offBudget: selectedAccount.offBudget
      });
      await getAccounts();
      toast({
        title: 'Account updated.',
        description: 'Account details have been successfully updated.',
        status: 'success',
        duration: 3000,
        isClosable: true
      });
      setIsModalOpen(false);
    } catch (err) {
      toast({
        title: 'Error updating account.',
        description: err.message,
        status: 'error',
        duration: 3000,
        isClosable: true
      });
    }
  };

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
    <Box>
      <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)}>
        <ModalOverlay />
        <ModalContent mx={{ base: '20px' }}>
          <ModalHeader>Edit Account</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            {selectedAccount && (
              <>
                <Input
                  placeholder="Account Name"
                  value={selectedAccount.name}
                  onChange={e =>
                    setSelectedAccount(prev => ({
                      ...prev,
                      name: e.target.value
                    }))
                  }
                  isInvalid={!!nameError}
                />
                {nameError && <Box color="red.500">{nameError}</Box>}

                <Select
                  placeholder="Select Category"
                  value={selectedAccount.category}
                  onChange={e =>
                    setSelectedAccount(prev => ({
                      ...prev,
                      category: e.target.value
                    }))
                  }
                  mt={3}
                  isInvalid={!!categoryError}
                >
                  {ACCOUNT_CATEGORIES.map(cat => (
                    <option key={cat} value={cat}>
                      {cat}
                    </option>
                  ))}
                </Select>
                {categoryError && <Box color="red.500">{categoryError}</Box>}

                <Checkbox
                  isChecked={selectedAccount.offBudget}
                  onChange={e =>
                    setSelectedAccount(prev => ({
                      ...prev,
                      offBudget: e.target.checked
                    }))
                  }
                  mt={3}
                >
                  Off Budget
                </Checkbox>
              </>
            )}
          </ModalBody>
          <ModalFooter>
            <Button colorScheme="blue" onClick={handleSave}>
              Save
            </Button>
            <Button onClick={() => setIsModalOpen(false)} ml={3}>
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>

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
              localAccounts.map((account, index) => (
                // eslint-disable-next-line react/no-array-index-key
                <Tr key={`${account.id}-${index}`}>
                  <Td>{account.name}</Td>
                  <Td>{account.category}</Td>
                  <Td>{account.offBudget ? 'Yes' : 'No'}</Td>
                  <Td>
                    <Button
                      colorScheme={primaryColor}
                      bg={primaryColor}
                      margin="2"
                      onClick={() => handleEdit(account)}
                    >
                      Edit
                    </Button>
                    <Button
                      colorScheme="red"
                      margin="2"
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
    </Box>
  );
};

export default Accounts;
