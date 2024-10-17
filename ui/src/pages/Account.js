import React from 'react';
import { useParams } from 'react-router-dom';
import { Box, Text } from '@chakra-ui/react';

const Home = () => {
  const { id } = useParams();
  return (
    <Box as="main" w="100%">
      <Text>Account - {id}</Text>
    </Box>
  );
};

export default Home;
