import React, { useState, useEffect } from 'react';
import {
  Box,
  Input,
  List,
  ListItem,
  ListIcon,
  Button,
  Text,
  Flex
} from '@chakra-ui/react';
import { AddIcon } from '@chakra-ui/icons';
import { usePayees } from '../context';

const AutocompletePayee = () => {
  const { payees, createPayee, error } = usePayees();
  const [localPayees, setLocalPayees] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [filteredPayees, setFilteredPayees] = useState([]);

  useEffect(() => {
    setLocalPayees(payees);
  }, [payees]);

  useEffect(() => {
    if (searchQuery) {
      setFilteredPayees(
        localPayees.filter(payee =>
          payee.name.toLowerCase().includes(searchQuery.toLowerCase())
        )
      );
    } else {
      setFilteredPayees([]);
    }
  }, [searchQuery, localPayees]);

  const handlePayeeSelect = payeeName => {
    setSearchQuery(payeeName);
    setFilteredPayees([]);
  };

  const handleAddNewPayee = async () => {
    if (searchQuery.trim()) {
      try {
        const newPayeeObj = {
          name: searchQuery.trim()
        };
        await createPayee(newPayeeObj);
        setSearchQuery('');
        setFilteredPayees([]);
      } catch (err) {
        // eslint-disable-next-line no-console
        console.error(err);
      }
    }
  };

  if (error) return <Box>Error loading payees: {error.message}</Box>;

  const payeeExists = localPayees.some(
    payee => payee.name.toLowerCase() === searchQuery.toLowerCase()
  );

  return (
    <Box overflowX="auto" borderRadius="md">
      <Box marginTop="1rem" bg="white" padding="1rem">
        <Input
          value={searchQuery}
          onChange={e => setSearchQuery(e.target.value)}
          placeholder="Search Payee"
          maxLength={50}
        />
        {filteredPayees.length > 0 && (
          <List
            spacing={1}
            marginTop="0.5rem"
            border="1px solid"
            borderColor="gray.300"
            borderRadius="md"
          >
            {filteredPayees.map(payee => (
              <ListItem
                key={payee.id}
                padding="0.5rem"
                cursor="pointer"
                onClick={() => handlePayeeSelect(payee.name)}
                _hover={{ bg: 'gray.100' }}
              >
                <Flex align="center">
                  <ListIcon as={AddIcon} boxSize={3} />{' '}
                  <Text>{payee.name}</Text>{' '}
                </Flex>
              </ListItem>
            ))}
          </List>
        )}
        {!payeeExists && searchQuery && (
          <Box marginTop="0.5rem">
            <Text color="gray.600">
              Payee not found. Would you like to add it?
            </Text>
            <Button marginTop="0.5rem" onClick={handleAddNewPayee}>
              Add `{searchQuery}`
            </Button>
          </Box>
        )}
      </Box>
    </Box>
  );
};

export default AutocompletePayee;
