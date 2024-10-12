import React from 'react';
import { VStack, Link } from '@chakra-ui/react';
import { Link as RouterLink, useLocation } from 'react-router-dom';

const NavigationMenu = ({ onClose }) => {
  const location = useLocation();

  const currentPath = location.pathname.replace('/', '');

  return (
    <VStack spacing={2} align="flex-start" width="100%">
      {['home', 'about', 'contact'].map(item => (
        <Link
          key={item}
          as={RouterLink}
          to={`/${item}`}
          onClick={() => {
            onClose();
          }}
          width="100%"
          padding="0.5rem 1rem"
          bg={currentPath === item ? 'gray.100' : 'transparent'}
          color="black"
          borderRadius="md"
          transition="background 0.2s"
          _hover={{ bg: 'gray.100' }}
        >
          {item.charAt(0).toUpperCase() + item.slice(1)}
        </Link>
      ))}
    </VStack>
  );
};

export default NavigationMenu;
