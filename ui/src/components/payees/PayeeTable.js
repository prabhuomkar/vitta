import React from 'react';
import { Table, Thead, Tbody, Tr, Th, Td } from '@chakra-ui/react';
import PayeeRow from './PayeeRow';
import NewPayeeRow from './NewPayeeRow';

const PayeeTable = ({
  payees,
  categories,
  newPayee,
  setNewPayee,
  handleAddPayee,
  handleDelete,
  handleFieldChange,
  handleSaveChanges,
  errors
}) => {
  return (
    <Table
      size="sm"
      variant="simple"
      bg="white"
      border="gray.300"
      borderRadius="md"
    >
      <Thead>
        <Tr>
          <Th padding="0.6rem">Payee Name</Th>
          <Th padding="0.6rem">Category</Th>
          <Th padding="0.6rem">Rules</Th>
          <Th padding="0.6rem">Add / Delete</Th>
        </Tr>
      </Thead>
      <Tbody>
        <NewPayeeRow
          newPayee={newPayee}
          categories={categories}
          setNewPayee={setNewPayee}
          handleAddPayee={handleAddPayee}
        />
        {payees.length > 0 ? (
          payees.map(payee => (
            <PayeeRow
              key={payee.id}
              payee={payee}
              categories={categories}
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
      </Tbody>
    </Table>
  );
};

export default PayeeTable;
