import React from 'react';
import { Box, Stack, Skeleton } from '@chakra-ui/react';

const Loading = () => (
  <Box bg="white" padding="1rem" borderRadius="md">
    <Stack>
      <Skeleton height="60px" />
      <Skeleton height="60px" />
    </Stack>
  </Box>
);

export default Loading;
