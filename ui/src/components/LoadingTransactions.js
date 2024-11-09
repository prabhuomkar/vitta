import React from 'react';
import {
  Box,
  Skeleton,
  Table,
  Thead,
  Tr,
  Th,
  Tbody,
  Td
} from '@chakra-ui/react';

const LoadingTransactions = () => (
  <Box bg="white" padding="" borderRadius="md" overflowX="auto">
    <Table variant="simple">
      <Thead>
        <Tr>
          <Th padding="0.6rem" width="18%">
            Name
          </Th>
          <Th padding="0.6rem" width="13%">
            Payee
          </Th>
          <Th padding="0.6rem" width="17%">
            Category
          </Th>
          <Th padding="0.6rem" width="10%">
            Credit
          </Th>
          <Th padding="0.6rem" width="10%">
            Debit
          </Th>
          <Th padding="0.6rem">Notes</Th>
          <Th padding="0.6rem">Cleared</Th>
          <Th padding="0.6rem">Actions</Th>
        </Tr>
      </Thead>
      <Tbody>
        {Array.from({ length: 4 }).map((x, rowIndex) => (
          // eslint-disable-next-line react/no-array-index-key
          <Tr key={rowIndex}>
            {Array.from({ length: 8 }).map((y, colIndex) => (
              // eslint-disable-next-line react/no-array-index-key
              <Td padding="0.6rem" key={colIndex}>
                <Skeleton height="28px" />
              </Td>
            ))}
          </Tr>
        ))}
      </Tbody>
    </Table>
  </Box>
);

export default LoadingTransactions;
