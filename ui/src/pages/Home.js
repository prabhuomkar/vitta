import React from 'react';
import { Box } from '@chakra-ui/react';
import Account from './Account';

const Home = () => {
  return (
    <Box as="main" w="100%" padding="1rem">
      {/* <Heading>Home</Heading> */}
      <Account />
    </Box>
  );
};

export default Home;
