/* eslint-disable no-nested-ternary */
import React, { useState, useContext } from 'react';
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
  useTheme,
  useToast
} from '@chakra-ui/react';
import { PlusSquareIcon } from '@chakra-ui/icons';
import { AccountsContext } from '../../context/AccountsContext';
import { NAV_ITEMS, ACCOUNT_CATEGORIES } from '../common/constants';

const NavigationMenu = ({ onClose }) => {
  const location = useLocation();
  const navigate = useNavigate();
  const currentPath = location.pathname;
  const { accounts, createAccount, loading } = useContext(AccountsContext);

  const theme = useTheme();
  const primaryColor = theme.colors.primary;
  const toast = useToast();

  // Form state
  const [accountName, setAccountName] = useState('');
  const [accountCategory, setAccountCategory] = useState('');
  const [offBudget, setOffBudget] = useState(false);
  const [formError, setFormError] = useState('');
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);

  const handleCreateAccount = async () => {
    if (!accountName || !accountCategory) {
      setFormError('Account name and category are required.');
      return;
    }

    try {
      await createAccount({
        name: accountName,
        category: accountCategory,
        offBudget
      });

      toast({
        title: 'Account added.',
        description: 'The account has been successfully added.',
        status: 'success',
        duration: 3000,
        isClosable: true
      });

      setAccountName('');
      setAccountCategory('');
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
            padding="0.5rem 1rem"
            bg={currentPath === `/${item}` ? 'gray.100' : 'transparent'}
            color="black"
            borderRadius="md"
            transition="background 0.2s"
            _hover={{ bg: 'gray.100' }}
          >
            {item.charAt(0).toUpperCase() + item.slice(1)}
          </Link>
        ))}
        <Divider />
        <Flex
          align="center"
          width="100%"
          padding="0.5rem 1rem"
          justify="space-between"
        >
          <Text>Accounts</Text>
          <Popover
            placement="auto-end"
            isOpen={isPopoverOpen}
            onClose={() => setIsPopoverOpen(false)}
          >
            <PopoverTrigger>
              <IconButton
                aria-label="Add account"
                icon={<PlusSquareIcon />}
                variant="outline"
                onClick={() => setIsPopoverOpen(true)}
              />
            </PopoverTrigger>
            <Portal>
              <PopoverContent>
                <PopoverArrow />
                <PopoverCloseButton position="absolute" top="8px" />
                <PopoverHeader>Add Account</PopoverHeader>
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
                    />
                  </FormControl>
                  <FormControl
                    id="account-category"
                    mb={4}
                    isInvalid={formError && !accountCategory}
                  >
                    <FormLabel>Account Category</FormLabel>
                    <Select
                      placeholder="Select category"
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
        <Divider />
        {loading ? (
          <Skeleton width="100%">
            <div>Loading...</div>
          </Skeleton>
        ) : accounts.length === 0 ? (
          <Text padding="0.5rem 1rem" color="gray.500">
            No accounts available
          </Text>
        ) : (
          accounts.map(account => (
            <Link
              key={account.id}
              as={RouterLink}
              to={`/account/${account.id}`}
              onClick={onClose}
              width="100%"
              padding="0.5rem 1rem"
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
      </>
    </VStack>
  );
};

export default NavigationMenu;
