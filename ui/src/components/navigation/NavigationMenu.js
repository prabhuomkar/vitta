import React from 'react';
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
  FormLabel,
  useTheme
} from '@chakra-ui/react';
import { PlusSquareIcon } from '@chakra-ui/icons';
import { Link as RouterLink, useLocation } from 'react-router-dom';

const NavigationMenu = ({ onClose }) => {
  const location = useLocation();
  const currentPath = location.pathname.replace('/', '');

  const theme = useTheme();
  const primaryColor = theme.colors.primary;

  return (
    <VStack spacing={2} align="flex-start" width="100%">
      <>
        {['home', 'about', 'contact'].map(item => (
          <Link
            key={item}
            as={RouterLink}
            to={`/${item}`}
            onClick={() => {
              onClose();
            }}
            width="100%"
            padding="0.5rem 1rem"
            bg={currentPath === item ? 'gray.100' : 'transparent'}
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
          <Popover placement="auto-end">
            <PopoverTrigger>
              <IconButton
                aria-label="Add account"
                icon={<PlusSquareIcon />}
                variant="outline"
              />
            </PopoverTrigger>
            <PopoverContent>
              <PopoverArrow />
              <PopoverCloseButton />
              <PopoverHeader>Add Account</PopoverHeader>
              <PopoverBody>
                <FormControl id="account-name" mb={4}>
                  <FormLabel>Account Name</FormLabel>
                  <Input placeholder="Enter account name" />
                </FormControl>
                <FormControl id="account-category" mb={4}>
                  <FormLabel>Account Category</FormLabel>
                  <Select placeholder="Select category">
                    <option value="savings">Savings</option>
                    <option value="credit">Credit</option>
                    <option value="investment">Investment</option>
                  </Select>
                </FormControl>
                <FormControl mb={4}>
                  <Checkbox>Off Budget</Checkbox>
                </FormControl>
                <Flex justify="flex-end" mb={1}>
                  <Button colorScheme={primaryColor} bgColor={primaryColor}>
                    Submit
                  </Button>
                </Flex>
              </PopoverBody>
            </PopoverContent>
          </Popover>
        </Flex>
        <Divider />
      </>
    </VStack>
  );
};

export default NavigationMenu;
