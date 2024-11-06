import React from 'react';
import { Tr, Td, Input, IconButton } from '@chakra-ui/react';
import { AddIcon } from '@chakra-ui/icons';

const NewPayeeRow = ({ newPayee, setNewPayee, handleAddPayee }) => {
  return (
    <Tr>
      <Td padding="0.6rem">
        <Input
          value={newPayee}
          onChange={e => setNewPayee(e.target.value)}
          placeholder="New Payee Name"
          maxLength={50}
          size="sm"
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
