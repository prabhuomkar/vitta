import React, { useState } from 'react';
import {
  Tr,
  Td,
  Text,
  Input,
  Select,
  Button,
  useTheme
} from '@chakra-ui/react';

const NewPayeeRow = ({
  payees,
  newPayee,
  categories,
  setNewPayee,
  handleAddPayee
}) => {
  const theme = useTheme();
  const primaryColor = theme.colors.primary;
  const [selectedCategory, setSelectedCategory] = useState('');

  const handleKeyDown = e => {
    if (e.key === 'Enter') {
      handleAddPayee(newPayee, selectedCategory || null);
      setSelectedCategory('');
    }
  };

  const handleAdd = () => {
    handleAddPayee(newPayee, selectedCategory || null);
    setSelectedCategory('');
  };

  return (
    <Tr>
      <Td padding="0.6rem">
        <Input
          value={newPayee}
          onChange={e => setNewPayee(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="New Payee Name"
          size="sm"
          maxLength={255}
        />
      </Td>
      <Td padding="0.6rem">
        <Select
          placeholder="Select Category"
          value={selectedCategory}
          onChange={e => setSelectedCategory(e.target.value)}
          size="sm"
        >
          {categories.map(category => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </Select>
      </Td>
      <Td padding="0.6rem">
        {payees.length === 0 && (
          <Text fontSize="xs" color="gray.500">
            You can add rules <br /> after adding a payee
          </Text>
        )}
      </Td>
      <Td padding="0.6rem">
        <Button
          aria-label="Add payee"
          onClick={handleAdd}
          size="sm"
          colorScheme={primaryColor}
          bg={primaryColor}
        >
          Add Payee
        </Button>
      </Td>
    </Tr>
  );
};

export default NewPayeeRow;
