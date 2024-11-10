import React from 'react';
import {
  Box,
  Text,
  Drawer,
  DrawerBody,
  DrawerHeader,
  DrawerOverlay,
  DrawerContent,
  DrawerCloseButton
} from '@chakra-ui/react';
import NavigationMenu from './NavigationMenu';

const NavigationDrawer = ({ isOpen, onClose }) => {
  return (
    <Drawer placement="left" onClose={onClose} isOpen={isOpen} size="full">
      <DrawerOverlay>
        <DrawerContent>
          <DrawerCloseButton />
          <DrawerHeader />
          <DrawerBody>
            <NavigationMenu onClose={onClose} />
          </DrawerBody>
        </DrawerContent>
        <Box as="footer" width="100%" color="gray.400" textAlign="center">
          <Text fontSize="xs">&copy; {new Date().getFullYear()} Vitta</Text>
        </Box>
      </DrawerOverlay>
    </Drawer>
  );
};

export default NavigationDrawer;
