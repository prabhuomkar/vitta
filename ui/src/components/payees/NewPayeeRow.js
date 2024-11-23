import React from 'react';
import { Tr, Td, Input, IconButton } from '@chakra-ui/react';
import { AddIcon } from '@chakra-ui/icons';

const NewPayeeRow = ({ newPayee, setNewPayee, handleAddPayee }) => {
  const handleKeyDown = e => {
    if (e.key === 'Enter') {
      handleAddPayee();
    }
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
        <IconButton
          aria-label="Add payee"
          icon={<AddIcon />}
          variant="outline"
          onClick={handleAddPayee}
          size="sm"
        />
      </Td>
    </Tr>
  );
};

export default NewPayeeRow;
