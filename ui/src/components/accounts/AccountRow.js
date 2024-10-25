import React, { useState } from 'react';
import { Tr, Td, Input, Select, Checkbox, IconButton } from '@chakra-ui/react';
import { DeleteIcon } from '@chakra-ui/icons';
import { ACCOUNT_CATEGORIES } from '../common/constants';

const AccountRow = ({ account, onDelete, onSave }) => {
  const [localAccount, setLocalAccount] = useState(account);

  const handleFieldChange = (field, value) => {
    setLocalAccount({ ...localAccount, [field]: value });
  };

  const handleBlur = () => {
    onSave(localAccount.id, localAccount);
  };

  return (
    <Tr>
      <Td padding="0.6rem">
        <Input
          value={localAccount.name || ''}
          onChange={e => handleFieldChange('name', e.target.value)}
          onBlur={handleBlur}
          placeholder="Account Name"
        />
      </Td>
      <Td padding="0.6rem">
        <Select
          value={localAccount.category || ''}
          onChange={e => handleFieldChange('category', e.target.value)}
          onBlur={handleBlur}
          placeholder="Select Category"
        >
          {ACCOUNT_CATEGORIES.map(cat => (
            <option key={cat} value={cat}>
              {cat}
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
        />
      </Td>
    </Tr>
  );
};

export default AccountRow;
