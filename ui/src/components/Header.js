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
import { useAuth } from '../context';

const username = process.env.REACT_APP_USERNAME;
const avatar = username ? username.charAt(0).toUpperCase() : '';

const Header = ({ onOpen, logoSrc }) => {
  const { logout } = useAuth();
  return (
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
          <span style={{ fontFamily: 'Libre Baskerville, serif' }}>Vitta</span>
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
            {avatar}
          </Box>
        </MenuButton>
        <MenuList>
          <MenuItem isDisabled>Settings</MenuItem>
          <MenuItem onClick={logout}>Logout</MenuItem>
        </MenuList>
      </Menu>
    </Flex>
  );
};

export default Header;
