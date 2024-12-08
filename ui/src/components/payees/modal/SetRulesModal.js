/* eslint-disable no-use-before-define */
import React, { useState } from 'react';
import {
  Box,
  Flex,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Button,
  Input,
  FormControl,
  FormLabel,
  useTheme
} from '@chakra-ui/react';

const SetRulesModal = ({ isOpen, onClose, onSave, initialValues = {} }) => {
  const theme = useTheme();
  const primaryColor = theme.colors.primary;
  const [includes, setIncludes] = useState(
    initialValues.includes?.join(', ') || ''
  );
  const [excludes, setExcludes] = useState(
    initialValues.excludes?.join(', ') || ''
  );
  const [startsWith, setStartsWith] = useState(
    initialValues.startsWith?.join(', ') || ''
  );
  const [endsWith, setEndsWith] = useState(
    initialValues.endsWith?.join(', ') || ''
  );

  const handleSave = () => {
    const includesArray = convertInputToArray(includes);
    const excludesArray = convertInputToArray(excludes);
    const startsWithArray = convertInputToArray(startsWith);
    const endsWithArray = convertInputToArray(endsWith);

    onSave({
      includes: includesArray,
      excludes: excludesArray,
      startsWith: startsWithArray,
      endsWith: endsWithArray
    });
    onClose();
  };

  function convertInputToArray(input) {
    if (typeof input === 'string') {
      return input.split(',').map(value => value.trim());
    }
    return [];
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Payee Name Rules</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <FormControl mb={4}>
            <FormLabel>Includes</FormLabel>
            <Input
              value={includes}
              onChange={e => setIncludes(e.target.value)}
              placeholder="Payee Name Includes"
            />
          </FormControl>

          <FormControl mb={4}>
            <FormLabel>Excludes</FormLabel>
            <Input
              value={excludes}
              onChange={e => setExcludes(e.target.value)}
              placeholder="Payee Name Excludes"
            />
          </FormControl>

          <FormControl mb={4}>
            <FormLabel>Starts With</FormLabel>
            <Input
              value={startsWith}
              onChange={e => setStartsWith(e.target.value)}
              placeholder="Payee Name Starts With"
            />
          </FormControl>

          <FormControl>
            <FormLabel>Ends With</FormLabel>
            <Input
              value={endsWith}
              onChange={e => setEndsWith(e.target.value)}
              placeholder="Payee Name Ends With"
            />
          </FormControl>
        </ModalBody>

        <ModalFooter>
          <Flex w="100%" alignItems="end" gap={4}>
            <Box
              flex="1"
              bg="gray.100"
              color="gray.500"
              p={2}
              borderRadius="md"
              fontSize="xs"
              display="flex"
              alignItems="center"
            >
              Enter comma-separated values if you want to filter by multiple
              criteria. For example, &ldquo;value1, value2, value3&ldquo;
            </Box>
            <Button
              colorScheme={primaryColor}
              bg={primaryColor}
              onClick={handleSave}
            >
              Save Rules
            </Button>
          </Flex>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default SetRulesModal;
