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

const LoadingBudgets = () => (
  <Box bg="white" borderRadius="md" overflowX="auto">
    <Table variant="simple">
      <Thead>
        <Tr>
          <Th padding="0.6rem" width="40%">
            Group / Category
          </Th>
          <Th padding="0.6rem" width="20%">
            Budgeted
          </Th>
          <Th padding="0.6rem" width="20%">
            Spent
          </Th>
          <Th padding="0.6rem" width="20%">
            Balance
          </Th>
        </Tr>
      </Thead>
      <Tbody>
        {Array.from({ length: 4 }).map((x, rowIndex) => (
          // eslint-disable-next-line react/no-array-index-key
          <Tr key={rowIndex}>
            {Array.from({ length: 4 }).map((y, colIndex) => (
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

export default LoadingBudgets;
