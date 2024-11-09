import React from 'react';
import { Button, Flex, Text, useTheme } from '@chakra-ui/react';

const Error = ({ message }) => {
  const theme = useTheme();
  const primaryColor = theme.colors.primary;

  // eslint-disable-next-line no-console
  console.log(message);

  return (
    <Flex direction="column" align="center" justify="center" textAlign="center">
      <Text mb={4}>Sorry, something went wrong.</Text>
      <Button
        onClick={() => window.location.reload()}
        colorScheme={primaryColor}
        bg={primaryColor}
        size="sm"
      >
        Refresh
      </Button>
    </Flex>
  );
};

export default Error;
