import React from 'react';
import { Tr, Td, Input, IconButton } from '@chakra-ui/react';
import { DeleteIcon } from '@chakra-ui/icons';

const PayeeRow = ({
  payee,
  handleDelete,
  handleFieldChange,
  handleSaveChanges,
  isFieldInvalid
}) => {
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
        <IconButton
          aria-label="Delete payee"
          icon={<DeleteIcon />}
          variant="outline"
          onClick={() => handleDelete(payee.id)}
          size="sm"
        />
      </Td>
    </Tr>
  );
};

export default PayeeRow;
