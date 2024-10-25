import React from 'react';
import { Table, Thead, Tbody, Tr, Th, Td } from '@chakra-ui/react';
import PayeeRow from './PayeeRow';
import NewPayeeRow from './NewPayeeRow';

const PayeeTable = ({
  payees,
  newPayee,
  setNewPayee,
  handleAddPayee,
  handleDelete,
  handleFieldChange,
  handleSaveChanges,
  errors
}) => {
  return (
    <Table variant="simple" bg="white" border="gray.300" borderRadius="md">
      <Thead>
        <Tr>
          <Th padding="0.6rem">Payee Name</Th>
          <Th padding="0.6rem">Actions</Th>
        </Tr>
      </Thead>
      <Tbody>
        {payees.length > 0 ? (
          payees.map(payee => (
            <PayeeRow
              key={payee.id}
              payee={payee}
              handleDelete={handleDelete}
              handleFieldChange={handleFieldChange}
              handleSaveChanges={handleSaveChanges}
              isFieldInvalid={errors[payee.id]?.name}
            />
          ))
        ) : (
          <Tr>
            <Td
              padding="0.6rem"
              colSpan={2}
              textAlign="center"
              color="gray.500"
            >
              No payees available
            </Td>
          </Tr>
        )}
        <NewPayeeRow
          newPayee={newPayee}
          setNewPayee={setNewPayee}
          handleAddPayee={handleAddPayee}
        />
      </Tbody>
    </Table>
  );
};

export default PayeeTable;
