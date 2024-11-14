/* eslint-disable no-nested-ternary */
import React, { useState } from 'react';
import { Link as RouterLink, useLocation, useNavigate } from 'react-router-dom';
import {
  VStack,
  Link,
  Divider,
  Text,
  IconButton,
  Flex,
  Popover,
  PopoverTrigger,
  PopoverContent,
  PopoverArrow,
  PopoverCloseButton,
  PopoverHeader,
  PopoverBody,
  Button,
  Input,
  FormControl,
  Select,
  Checkbox,
  Portal,
  Skeleton,
  FormLabel,
  Box,
  useTheme,
  useToast
} from '@chakra-ui/react';
import { PlusSquareIcon } from '@chakra-ui/icons';
import { useAccounts } from '../../context';
import { NAV_ITEMS, ACCOUNT_CATEGORIES, ADAPTERS } from '../common/constants';

const NavigationMenu = ({ onClose }) => {
  const location = useLocation();
  const navigate = useNavigate();
  const currentPath = location.pathname;
  const { accounts, createAccount, loading } = useAccounts();

  const theme = useTheme();
  const primaryColor = theme.colors.primary;
  const toast = useToast();

  // Form state
  const [accountName, setAccountName] = useState('');
  const [accountCategory, setAccountCategory] = useState('');
  const [accountAdapter, setAccountAdapter] = useState('');
  const [offBudget, setOffBudget] = useState(false);
  const [formError, setFormError] = useState('');
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);

  const handleCreateAccount = async () => {
    if (!accountName || !accountCategory || !accountAdapter) {
      setFormError('Account name, category and adapter are required.');
      return;
    }

    // function to check duplicate account name
    const duplicateAccount = accounts.some(
      account => account.name === accountName.trim()
    );

    if (duplicateAccount) {
      toast({
        title: 'Duplicate Account',
        description: `An account with the name "${accountName}" already exists.`,
        status: 'warning',
        duration: 1500,
        isClosable: true
      });
      return;
    }

    try {
      await createAccount({
        name: accountName,
        category: accountCategory,
        adapter: accountAdapter,
        offBudget
      });

      toast({
        title: 'Account added.',
        description: 'The account has been successfully added.',
        status: 'success',
        duration: 1500,
        isClosable: true
      });

      setAccountName('');
      setAccountCategory('');
      setAccountAdapter('');
      setOffBudget(false);
      setFormError('');
      setIsPopoverOpen(false);
    } catch (err) {
      setFormError('Failed to add account. Please try again.');
    }
  };

  return (
    <VStack spacing={2} align="flex-start" width="100%">
      <>
        {NAV_ITEMS.map(item => (
          <Link
            key={item}
            as={RouterLink}
            to={`/${item}`}
            onClick={() => {
              onClose();
            }}
            width="100%"
            padding="0.4rem 0.8rem"
            bg={currentPath === `/${item}` ? 'gray.100' : 'transparent'}
            color="black"
            borderRadius="md"
            transition="background 0.2s"
            _hover={{ bg: 'gray.100' }}
          >
            {item.charAt(0).toUpperCase() + item.slice(1)}
          </Link>
        ))}
        <Flex
          align="center"
          width="100%"
          paddingLeft="0.8rem"
          justify="space-between"
        >
          <Text fontSize="sm" as="b" color="gray.600">
            ACCOUNTS
          </Text>
          <Popover
            placement="auto-end"
            isOpen={isPopoverOpen}
            onClose={() => setIsPopoverOpen(false)}
          >
            <PopoverTrigger>
              <IconButton
                aria-label="Add account"
                size="sm"
                icon={<PlusSquareIcon />}
                fontSize="16px"
                variant="outline"
                onClick={() => setIsPopoverOpen(true)}
              />
            </PopoverTrigger>
            <Portal>
              <PopoverContent>
                <PopoverArrow />
                <PopoverCloseButton position="absolute" top="8px" />
                <PopoverHeader>ADD ACCOUNT</PopoverHeader>
                <PopoverBody>
                  <FormControl
                    id="account-name"
                    mb={4}
                    isInvalid={formError && !accountName}
                  >
                    <FormLabel>Account Name</FormLabel>
                    <Input
                      placeholder="Enter account name"
                      value={accountName}
                      onChange={e => setAccountName(e.target.value)}
                      maxLength={255}
                    />
                  </FormControl>
                  <FormControl
                    id="account-category"
                    mb={4}
                    isInvalid={formError && !accountCategory}
                  >
                    <FormLabel>Account Category</FormLabel>
                    <Select
                      placeholder="Select Category"
                      value={accountCategory}
                      onChange={e => setAccountCategory(e.target.value)}
                    >
                      {ACCOUNT_CATEGORIES.map(category => (
                        <option key={category} value={category}>
                          {category}
                        </option>
                      ))}
                    </Select>
                  </FormControl>
                  <FormControl
                    id="account-adapter"
                    mb={4}
                    isInvalid={formError && !accountAdapter}
                  >
                    <FormLabel>Adapter</FormLabel>
                    <Select
                      placeholder="Select Adapter"
                      value={accountAdapter}
                      onChange={e => setAccountAdapter(e.target.value)}
                    >
                      {ADAPTERS.map(adapter => (
                        <option key={adapter} value={adapter}>
                          {adapter}
                        </option>
                      ))}
                    </Select>
                  </FormControl>
                  <FormControl mb={4}>
                    <Checkbox
                      isChecked={offBudget}
                      onChange={e => setOffBudget(e.target.checked)}
                    >
                      Off Budget
                    </Checkbox>
                  </FormControl>
                  {formError && (
                    <Text color="red.500" mb={4}>
                      {formError}
                    </Text>
                  )}
                  <Flex justify="flex-end" mb={1}>
                    <Button
                      colorScheme={primaryColor}
                      bgColor={primaryColor}
                      onClick={handleCreateAccount}
                    >
                      Submit
                    </Button>
                  </Flex>
                </PopoverBody>
              </PopoverContent>
            </Portal>
          </Popover>
        </Flex>
        {loading ? (
          <Skeleton width="100%">
            <div>Loading...</div>
          </Skeleton>
        ) : accounts.length === 0 ? (
          <Flex
            textAlign="center"
            justifyContent="center"
            alignItems="center"
            width="100%"
            flexDirection="column"
          >
            <Text fontSize="sm" padding="0.4rem 0.8rem" color="gray.500">
              No accounts available
            </Text>
            <Button
              onClick={() => setIsPopoverOpen(true)}
              colorScheme={primaryColor}
              bg={primaryColor}
              size="xs"
            >
              Add account
            </Button>
          </Flex>
        ) : (
          accounts.map(account => (
            <Link
              key={account.id}
              as={RouterLink}
              to={`/account/${account.id}`}
              onClick={onClose}
              width="100%"
              padding="0.4rem 0.8rem"
              bg={
                currentPath === `/account/${account.id}`
                  ? 'gray.100'
                  : 'transparent'
              }
              color="black"
              borderRadius="md"
              transition="background 0.2s"
              _hover={{ bg: 'gray.100' }}
            >
              {account.name}
            </Link>
          ))
        )}
        <Divider />
        {accounts.length > 0 && (
          <Button
            width="100%"
            colorScheme={primaryColor}
            bgColor={primaryColor}
            onClick={() => {
              navigate('/accounts');
              onClose();
            }}
          >
            Manage Accounts
          </Button>
        )}
      </>
      <Box as="footer" width="100%" color="gray.400" textAlign="center">
        <Text fontSize="xs">&copy; {new Date().getFullYear()} Vitta</Text>
      </Box>
    </VStack>
  );
};

export default NavigationMenu;
