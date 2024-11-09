import React, { useState } from 'react';
import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Button,
  FormControl,
  FormLabel,
  Input,
  Textarea,
  Box,
  useToast,
  useTheme
} from '@chakra-ui/react';
import { useTransactions } from '../../context';

const AddTransactionModal = ({ isOpen, onClose, getTransactions }) => {
  const theme = useTheme();
  const primaryColor = theme.colors.primary;
  const [name, setName] = useState('');
  const [credit, setCredit] = useState('');
  const [debit, setDebit] = useState('');
  const [notes, setNotes] = useState('');
  const { createTransaction } = useTransactions();

  const [errors, setErrors] = useState({});
  const toast = useToast();

  const handleSubmit = async () => {
    const newErrors = {};
    if (!name) newErrors.name = 'Name is required';

    setErrors(newErrors);

    if (Object.keys(newErrors).length > 0) return;

    const transactionData = {
      name,
      credit: parseFloat(credit) || 0,
      debit: parseFloat(debit) || 0,
      notes
    };

    try {
      await createTransaction(transactionData);
      await getTransactions();

      toast({
        title: 'Transaction added.',
        description: 'The transaction was added successfully.',
        status: 'success',
        duration: 1500,
        isClosable: true
      });

      setName('');
      setCredit('');
      setDebit('');
      setNotes('');
      setErrors({});
      onClose();
    } catch (error) {
      toast({
        title: 'Error adding transaction.',
        description:
          'There was an error while adding the transaction. Please try again.',
        status: 'error',
        duration: 1500,
        isClosable: true
      });
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent mx={{ base: '4', md: '0' }}>
        <ModalHeader>ADD TRANSACTION</ModalHeader>
        <ModalCloseButton position="absolute" right="8px" />
        <ModalBody>
          <FormControl isInvalid={errors.name}>
            <FormLabel>Name</FormLabel>
            <Input
              value={name}
              onChange={e => setName(e.target.value)}
              placeholder="Transaction Name"
              maxLength={255}
            />
            {errors.name && <Box color="red.500">{errors.name}</Box>}
          </FormControl>
          <FormControl mt={4} isInvalid={errors.credit}>
            <FormLabel>Credit</FormLabel>
            <Input
              value={credit || 0}
              onChange={e => setCredit(e.target.value)}
              placeholder="Credit Amount"
              type="number"
              step="0.01"
            />
            {errors.credit && <Box color="red.500">{errors.credit}</Box>}
          </FormControl>
          <FormControl mt={4} isInvalid={errors.debit}>
            <FormLabel>Debit</FormLabel>
            <Input
              value={debit || 0}
              onChange={e => setDebit(e.target.value)}
              placeholder="Debit Amount"
              type="number"
              step="0.01"
            />
            {errors.debit && <Box color="red.500">{errors.debit}</Box>}
          </FormControl>
          <FormControl mt={4}>
            <FormLabel>Notes</FormLabel>
            <Textarea
              value={notes}
              onChange={e => setNotes(e.target.value)}
              placeholder="Additional Notes"
              maxLength={512}
            />
          </FormControl>
        </ModalBody>

        <ModalFooter>
          <Button
            colorScheme={primaryColor}
            bg={primaryColor}
            onClick={handleSubmit}
          >
            Add Transaction
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default AddTransactionModal;
