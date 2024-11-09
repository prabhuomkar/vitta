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
  FormErrorMessage,
  Textarea,
  Button
} from '@chakra-ui/react';

const AddGroupModal = ({
  isOpen,
  onClose,
  formData,
  groupNameError,
  handleChange,
  handleAddGroupSubmit,
  primaryColor
}) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent mx={{ base: '4', md: '0' }}>
        <ModalHeader>Add Group</ModalHeader>
        <ModalCloseButton position="absolute" right="8px" />
        <ModalBody>
          <form onSubmit={handleAddGroupSubmit}>
            <FormControl mb={4} isInvalid={!!groupNameError}>
              <FormLabel>Group Name</FormLabel>
              <Input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                placeholder="Enter Group Name"
                maxLength={255}
              />
              {groupNameError && (
                <FormErrorMessage>{groupNameError}</FormErrorMessage>
              )}
            </FormControl>
            <FormControl mb={4}>
              <FormLabel>Notes</FormLabel>
              <Textarea
                name="notes"
                value={formData.notes}
                onChange={handleChange}
                placeholder="Enter Notes"
                maxLength={512}
              />
            </FormControl>
            <ModalFooter margin="auto" mb={2} p={0}>
              <Button
                colorScheme={primaryColor}
                bg={primaryColor}
                type="submit"
              >
                Add Group
              </Button>
            </ModalFooter>
          </form>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};

export default AddGroupModal;
