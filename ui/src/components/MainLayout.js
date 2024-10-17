import React, { useState } from 'react';
import { Box, Flex, useDisclosure } from '@chakra-ui/react';
import Header from './Header';
import { NavigationDrawer, NavigationMenu } from './navigation';

const MainLayout = ({ children }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [activeLink, setActiveLink] = useState('');

  const handleLinkClick = link => {
    setActiveLink(link);
    onClose();
  };

  return (
    <Box bg="gray.100" minH="100vh">
      <Header onOpen={onOpen} logoSrc="/logo.png" />
      <Flex>
        <Box
          as="nav"
          display={{ base: 'none', md: 'block' }}
          position="sticky"
          top="64px"
          w="275px"
          h="calc(100vh - 64px)"
          bg="white"
          padding="1rem"
          borderRight="1px solid"
          borderColor="gray.200"
          zIndex="9"
          transform="translateZ(0)"
          overflowY="scroll"
        >
          <NavigationMenu onClose={onClose} />
        </Box>
        <Box
          as="main"
          w="100%"
          minWidth={['100%', '100%', '50%']}
          padding="1rem"
        >
          {children}
        </Box>
      </Flex>
      <NavigationDrawer
        isOpen={isOpen}
        onClose={onClose}
        activeLink={activeLink}
        onLinkClick={handleLinkClick}
      />
    </Box>
  );
};

export default MainLayout;
