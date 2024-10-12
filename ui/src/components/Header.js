import {
  Box,
  IconButton,
  Flex,
  Heading,
  Image,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  Button
} from '@chakra-ui/react';
import { HamburgerIcon } from '@chakra-ui/icons';

const Header = ({ onOpen, logoSrc }) => (
  <Flex
    as="header"
    align="center"
    justify="space-between"
    padding="1rem"
    bg="white"
    color="black"
    borderBottom="1px solid"
    borderColor="gray.200"
    position="sticky"
    top="0"
    zIndex="10"
    w="100%"
    h="64px"
    transform="translateZ(0)"
  >
    <IconButton
      display={{ base: 'block', md: 'none' }}
      icon={<HamburgerIcon color="black" />}
      onClick={onOpen}
      fontSize="24px"
      variant="unstyled"
      aria-label="Open Menu"
    />
    <Flex align="center">
      <Image
        src={logoSrc}
        alt="Logo"
        boxSize="40px"
        display={{ base: 'none', md: 'block' }}
      />
      <Heading size="md" ml="1rem" display={{ base: 'none', md: 'block' }}>
        Vitta
      </Heading>
    </Flex>
    <Image
      src={logoSrc}
      alt="Logo"
      boxSize="40px"
      mx="auto"
      display={{ base: 'block', md: 'none' }}
    />
    <Menu>
      <MenuButton
        as={Button}
        rounded="full"
        variant="link"
        cursor="pointer"
        minW={0}
        _hover={{ border: 'none', bg: 'transparent' }}
        _focus={{ boxShadow: 'none', outline: 'none' }}
      >
        <Box
          width="40px"
          height="40px"
          borderRadius="full"
          bg="gray.300"
          display="flex"
          alignItems="center"
          justifyContent="center"
        >
          AP
        </Box>
      </MenuButton>
      <MenuList>
        <MenuItem>Settings</MenuItem>
        <MenuItem>Logout</MenuItem>
      </MenuList>
    </Menu>
  </Flex>
);

export default Header;
