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

const AddCategoryModal = ({
  isOpen,
  onClose,
  categoryFormData,
  categoryNameError,
  handleCategoryChange,
  handleAddCategorySubmit,
  primaryColor
}) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent mx={{ base: '4', md: '0' }}>
        <ModalHeader>Add Category</ModalHeader>
        <ModalCloseButton position="absolute" right="8px" />
        <ModalBody>
          <form onSubmit={handleAddCategorySubmit}>
            <FormControl mb={4} isInvalid={!!categoryNameError}>
              <FormLabel>Category Name</FormLabel>
              <Input
                type="text"
                name="name"
                value={categoryFormData.name}
                onChange={handleCategoryChange}
                placeholder="Enter Category Name"
                maxLength={255}
              />
              {categoryNameError && (
                <FormErrorMessage>{categoryNameError}</FormErrorMessage>
              )}
            </FormControl>
            <FormControl mb={4}>
              <FormLabel>Notes</FormLabel>
              <Textarea
                type="text"
                name="notes"
                value={categoryFormData.notes}
                onChange={handleCategoryChange}
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
                Add Category
              </Button>
            </ModalFooter>
          </form>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};

export default AddCategoryModal;
