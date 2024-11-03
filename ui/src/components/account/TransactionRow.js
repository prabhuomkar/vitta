import React from 'react';
import { Tr, Td, Input, Select, Checkbox, IconButton } from '@chakra-ui/react';
import { DeleteIcon } from '@chakra-ui/icons';

const TransactionRow = ({
  transaction,
  index,
  handleInputChange,
  handleSaveChanges,
  validationErrors,
  payees,
  categories,
  handleCheckboxChange,
  handleDelete
}) => (
  <Tr>
    <Td padding="0.5rem">
      <Input
        value={transaction.name || ''}
        onChange={e => handleInputChange(index, 'name', e.target.value)}
        onBlur={() => handleSaveChanges(index)}
        borderColor={validationErrors[index]?.name ? 'red.500' : undefined}
      />
    </Td>
    <Td padding="0.5rem">
      <Select
        placeholder="Select Payee"
        value={transaction.payeeId || ''}
        onChange={e => handleInputChange(index, 'payeeId', e.target.value)}
        onBlur={() => handleSaveChanges(index)}
        borderColor={validationErrors[index]?.payeeId ? 'red.500' : undefined}
      >
        {payees.map(payee => (
          <option key={payee.id} value={payee.id}>
            {payee.name}
          </option>
        ))}
      </Select>
    </Td>
    <Td padding="0.5rem">
      <Select
        placeholder="Select Category"
        value={transaction.categoryId || ''}
        onChange={e => handleInputChange(index, 'categoryId', e.target.value)}
        onBlur={() => handleSaveChanges(index)}
        borderColor={
          validationErrors[index]?.categoryId ? 'red.500' : undefined
        }
      >
        {categories.map(category => (
          <option key={category.id} value={category.id}>
            {category.name}
          </option>
        ))}
      </Select>
    </Td>
    <Td padding="0.5rem" width="10%">
      <Input
        type="number"
        value={transaction.credit || 0}
        onChange={e =>
          handleInputChange(index, 'credit', parseFloat(e.target.value) || 0)
        }
        onBlur={() => handleSaveChanges(index)}
      />
    </Td>
    <Td padding="0.5rem" width="10%">
      <Input
        type="number"
        value={transaction.debit || 0}
        onChange={e =>
          handleInputChange(index, 'debit', parseFloat(e.target.value) || 0)
        }
        onBlur={() => handleSaveChanges(index)}
      />
    </Td>
    <Td padding="0.5rem">
      <Input
        value={transaction.notes || ''}
        onChange={e => handleInputChange(index, 'notes', e.target.value)}
        onBlur={() => handleSaveChanges(index)}
      />
    </Td>
    <Td padding="0.5rem">
      <Checkbox
        isChecked={!!transaction.clearedAt}
        onChange={() => handleCheckboxChange(index)}
        onBlur={() => handleSaveChanges(index)}
      />
    </Td>
    <Td padding="0.6rem">
      <IconButton
        aria-label="Delete account"
        icon={<DeleteIcon />}
        variant="outline"
        onClick={() => handleDelete(transaction.id)}
      />
    </Td>
  </Tr>
);

export default TransactionRow;
