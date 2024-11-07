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

const EditCategoryModal = ({
  isOpen,
  onClose,
  editCategoryFormData,
  categoryNameError,
  handleEditCategoryChange,
  handleEditCategorySubmit,
  handleDeleteCategory,
  primaryColor
}) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent mx={{ base: '4', md: '0' }}>
        <ModalHeader>Edit Category</ModalHeader>
        <ModalCloseButton position="absolute" right="8px" />
        <ModalBody>
          <form onSubmit={handleEditCategorySubmit}>
            <FormControl mb={4} isInvalid={!!categoryNameError}>
              <FormLabel>Category Name</FormLabel>
              <Input
                type="text"
                name="name"
                value={editCategoryFormData.name}
                onChange={handleEditCategoryChange}
                placeholder="Enter Category Name"
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
                value={editCategoryFormData.notes}
                onChange={handleEditCategoryChange}
                placeholder="Enter Notes"
              />
            </FormControl>
            <ModalFooter margin="auto" mb={2} p={0}>
              <Button
                colorScheme={primaryColor}
                bg={primaryColor}
                type="submit"
              >
                Update
              </Button>
            </ModalFooter>
          </form>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};

export default EditCategoryModal;
