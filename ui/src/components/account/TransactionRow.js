import React, { useState } from 'react';
import { Tr, Td, Input, Select, Checkbox, IconButton } from '@chakra-ui/react';
import { DeleteIcon } from '@chakra-ui/icons';
import isEqual from 'lodash.isequal';

const TransactionRow = ({
  transaction,
  payees,
  categories,
  updateTransaction,
  deleteTransaction,
  toast
}) => {
  const [localTransaction, setLocalTransaction] = useState(transaction);
  const [validationErrors, setValidationErrors] = useState({});

  const handleInputChange = (field, value) => {
    setLocalTransaction(prev => ({
      ...prev,
      [field]: value || null
    }));

    if (field === 'name' && value) {
      setValidationErrors(prevErrors => ({ ...prevErrors, name: false }));
    }
  };

  const handleCheckboxChange = () => {
    setLocalTransaction(prev => ({
      ...prev,
      clearedAt: prev.clearedAt ? null : new Date().toISOString()
    }));
  };

  const handleSaveChanges = async () => {
    const errors = {};
    if (!localTransaction.name) {
      errors.name = true;
    }

    if (Object.keys(errors).length > 0) {
      setValidationErrors(errors);
      return;
    }

    const hasChanges = !isEqual(localTransaction, transaction);
    if (hasChanges) {
      await updateTransaction(localTransaction.id, localTransaction);
    }
  };

  const handleDelete = async () => {
    await deleteTransaction(localTransaction.id);
    toast({
      title: 'Transaction deleted.',
      description: `Transaction has been successfully deleted.`,
      status: 'success',
      duration: 1500,
      isClosable: true
    });
  };

  return (
    <Tr>
      <Td padding="0.6rem">
        <Input
          value={localTransaction.name || ''}
          onChange={e => handleInputChange('name', e.target.value)}
          onBlur={handleSaveChanges}
          borderColor={validationErrors.name ? 'red.500' : undefined}
          size="sm"
          maxLength={255}
        />
      </Td>
      <Td padding="0.6rem">
        <Select
          placeholder="Select Payee"
          value={localTransaction.payeeId || ''}
          onChange={e => handleInputChange('payeeId', e.target.value)}
          onBlur={handleSaveChanges}
          size="sm"
        >
          {payees.map(payee => (
            <option key={payee.id} value={payee.id}>
              {payee.name}
            </option>
          ))}
        </Select>
      </Td>
      <Td padding="0.6rem">
        <Select
          placeholder="Select Category"
          value={localTransaction.categoryId || ''}
          onChange={e => handleInputChange('categoryId', e.target.value)}
          onBlur={handleSaveChanges}
          size="sm"
        >
          {categories.map(category => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </Select>
      </Td>
      <Td padding="0.6rem" width="10%">
        <Input
          type="number"
          value={localTransaction.credit || 0}
          onChange={e =>
            handleInputChange('credit', parseFloat(e.target.value) || 0)
          }
          onBlur={handleSaveChanges}
          size="sm"
        />
      </Td>
      <Td padding="0.6rem" width="10%">
        <Input
          type="number"
          value={localTransaction.debit || 0}
          onChange={e =>
            handleInputChange('debit', parseFloat(e.target.value) || 0)
          }
          onBlur={handleSaveChanges}
          size="sm"
        />
      </Td>
      <Td padding="0.6rem">
        <Input
          value={localTransaction.notes || ''}
          onChange={e => handleInputChange('notes', e.target.value)}
          onBlur={handleSaveChanges}
          size="sm"
          maxLength={512}
        />
      </Td>
      <Td padding="0.6rem">
        <Checkbox
          isChecked={!!localTransaction.clearedAt}
          onChange={handleCheckboxChange}
          onBlur={handleSaveChanges}
        />
      </Td>
      <Td padding="0.6rem">
        <IconButton
          aria-label="Delete transaction"
          icon={<DeleteIcon />}
          variant="outline"
          onClick={handleDelete}
          size="sm"
        />
      </Td>
    </Tr>
  );
};

export default TransactionRow;
