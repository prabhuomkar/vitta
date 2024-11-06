import React from 'react';
import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  FormControl,
  FormLabel,
  Input,
  Button
} from '@chakra-ui/react';

const SetBudgetModal = ({
  isOpen,
  onClose,
  budgetFormData,
  setBudgetFormData,
  handleSetBudget,
  primaryColor
}) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent mx={{ base: '4', md: '0' }}>
        <ModalHeader>Add Budget</ModalHeader>
        <ModalCloseButton position="absolute" right="8px" />
        <ModalBody>
          <form onSubmit={handleSetBudget}>
            <FormControl mb={4}>
              <FormLabel>Budget</FormLabel>
              <Input
                type="number"
                name="budgeted"
                value={budgetFormData.budgeted}
                onChange={e => setBudgetFormData({ budgeted: e.target.value })}
                placeholder="Enter Budget Amount"
              />
            </FormControl>
            <ModalFooter margin="auto" mb={2} p={0}>
              <Button
                colorScheme={primaryColor}
                bg={primaryColor}
                type="submit"
              >
                Set Budget
              </Button>
            </ModalFooter>
          </form>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};

export default SetBudgetModal;
