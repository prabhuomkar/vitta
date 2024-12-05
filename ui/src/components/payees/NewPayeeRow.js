import React, { useState } from 'react';
import { Tr, Td, Input, Select, Tooltip, IconButton } from '@chakra-ui/react';
import { AddIcon } from '@chakra-ui/icons';
import { FaListUl } from 'react-icons/fa';

const NewPayeeRow = ({ newPayee, categories, setNewPayee, handleAddPayee }) => {
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
        <Tooltip label="Set Rules" hasArrow openDelay={500}>
          <IconButton icon={<FaListUl />} variant="outline" size="sm" />
        </Tooltip>
      </Td>
      <Td padding="0.6rem">
        <Tooltip label="Add Payee" hasArrow openDelay={500}>
          <IconButton
            aria-label="Add payee"
            icon={<AddIcon />}
            variant="outline"
            onClick={handleAdd}
            size="sm"
          />
        </Tooltip>
      </Td>
    </Tr>
  );
};

export default NewPayeeRow;
