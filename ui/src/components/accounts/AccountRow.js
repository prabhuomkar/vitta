import React, { useState } from 'react';
import { Tr, Td, Input, Select, Checkbox, IconButton } from '@chakra-ui/react';
import { DeleteIcon } from '@chakra-ui/icons';

const AccountRow = ({ account, adapters, onDelete, onSave }) => {
  const [localAccount, setLocalAccount] = useState(account);
  const [isNameInvalid, setIsNameInvalid] = useState(false);

  const handleFieldChange = (field, value) => {
    if (field === 'name') {
      setIsNameInvalid(!value.trim());
    }
    setLocalAccount({ ...localAccount, [field]: value });
  };

  const handleBlur = () => {
    if (!localAccount.name.trim()) {
      setIsNameInvalid(true);
      return;
    }
    onSave(localAccount.id, localAccount);
  };

  const ACCOUNT_ADAPTERS = [...new Set(adapters.map(item => item.name))];
  const ACCOUNT_CATEGORIES = [...new Set(adapters.map(item => item.category))];

  return (
    <Tr>
      <Td padding="0.6rem">
        <Input
          value={localAccount.name || ''}
          onChange={e => handleFieldChange('name', e.target.value)}
          onBlur={handleBlur}
          placeholder="Account Name"
          size="sm"
          maxLength={255}
          isInvalid={isNameInvalid}
          errorBorderColor="red.500"
        />
      </Td>
      <Td padding="0.6rem">
        <Select
          value={localAccount.category || ''}
          onChange={e => handleFieldChange('category', e.target.value)}
          onBlur={handleBlur}
          placeholder="Select Category"
          size="sm"
        >
          {ACCOUNT_CATEGORIES.map(cat => (
            <option key={cat} value={cat}>
              {cat}
            </option>
          ))}
        </Select>
      </Td>
      <Td padding="0.6rem">
        <Select
          value={localAccount.adapter || ''}
          onChange={e => handleFieldChange('adapter', e.target.value)}
          onBlur={handleBlur}
          placeholder="Select Adapter"
          size="sm"
        >
          {ACCOUNT_ADAPTERS.map(adapter => (
            <option key={adapter} value={adapter}>
              {adapter}
            </option>
          ))}
        </Select>
      </Td>
      <Td padding="0.6rem">
        <Checkbox
          isChecked={localAccount.offBudget || false}
          onChange={e => handleFieldChange('offBudget', e.target.checked)}
          onBlur={handleBlur}
        />
      </Td>
      <Td padding="0.6rem">
        <IconButton
          aria-label="Delete account"
          icon={<DeleteIcon />}
          variant="outline"
          onClick={() => onDelete(localAccount.id)}
          size="sm"
        />
      </Td>
    </Tr>
  );
};

export default AccountRow;
