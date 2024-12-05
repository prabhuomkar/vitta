import React from 'react';
import { Tr, Td, Input, Select, Tooltip, IconButton } from '@chakra-ui/react';
import { DeleteIcon } from '@chakra-ui/icons';
import { FaListUl } from 'react-icons/fa';

const PayeeRow = ({
  payee,
  categories,
  handleDelete,
  handleFieldChange,
  handleSaveChanges,
  isFieldInvalid
}) => {
  const handleCategoryChange = e => {
    const newCategoryId = e.target.value;
    handleFieldChange(payee.id, 'autoCategoryId', newCategoryId);
  };

  const handleCategoryBlur = () => {
    handleSaveChanges(payee.id, {
      ...payee,
      autoCategoryId: payee.autoCategoryId
    });
  };

  return (
    <Tr>
      <Td padding="0.6rem">
        <Input
          value={payee.name || ''}
          onChange={e => handleFieldChange(payee.id, 'name', e.target.value)}
          onBlur={() =>
            handleSaveChanges(payee.id, { ...payee, name: payee.name })
          }
          placeholder="Payee Name"
          isInvalid={isFieldInvalid}
          errorBorderColor="red.500"
          size="sm"
          maxLength={255}
        />
      </Td>
      <Td padding="0.6rem">
        <Select
          value={payee.autoCategoryId || ''}
          onChange={handleCategoryChange}
          onBlur={handleCategoryBlur}
          placeholder="Select Category"
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
        <Tooltip label="Delete Payee" hasArrow openDelay={500}>
          <IconButton
            aria-label="Delete payee"
            icon={<DeleteIcon />}
            variant="outline"
            onClick={() => handleDelete(payee.id)}
            size="sm"
          />
        </Tooltip>
      </Td>
    </Tr>
  );
};

export default PayeeRow;
