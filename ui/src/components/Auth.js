import React, { useState } from 'react';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Stack,
  useTheme,
  Image
} from '@chakra-ui/react';
import { useAuth } from '../context/AuthContext';

const Auth = () => {
  const theme = useTheme();
  const primaryColor = theme.colors.primary;
  const { login } = useAuth();

  const [formData, setFormData] = useState({
    username: '',
    password: ''
  });

  const [errors, setErrors] = useState({
    username: '',
    password: '',
    login: ''
  });

  const handleChange = e => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });

    if (name === 'username') {
      setErrors({ ...errors, username: '' });
    } else if (name === 'password') {
      setErrors({ ...errors, password: '' });
    }
  };

  const handleSubmit = () => {
    let hasError = false;

    if (!formData.username) {
      setErrors(prev => ({ ...prev, username: 'Username is required' }));
      hasError = true;
    }

    if (!formData.password) {
      setErrors(prev => ({ ...prev, password: 'Password is required' }));
      hasError = true;
    }

    if (!hasError) {
      try {
        login(formData.username, formData.password);
      } catch (error) {
        setErrors({
          username: '',
          password: '',
          login: error.message
        });
      }
    }
  };

  return (
    <Box
      minH="100vh"
      display="flex"
      alignItems="center"
      justifyContent="center"
      px={4}
      bg="gray.100"
    >
      <Box
        w={{ base: 'full', md: '400px' }}
        p={6}
        bg="white"
        borderRadius="md"
        border="1px solid"
        borderColor="gray.200"
      >
        <Stack spacing={4} align="center">
          <Image src="/logo.png" alt="Logo" boxSize="64px" mb={2} />
          <FormControl id="username">
            <FormLabel>Username</FormLabel>
            <Input
              type="text"
              name="username"
              placeholder="Enter your username"
              onChange={handleChange}
              borderColor={errors.username ? 'red.300' : 'gray.300'}
              focusBorderColor={primaryColor}
            />
            {errors.username && <Box color="red.500">{errors.username}</Box>}
          </FormControl>
          <FormControl id="password">
            <FormLabel>Password</FormLabel>
            <Input
              type="password"
              name="password"
              placeholder="Enter your password"
              onChange={handleChange}
              borderColor={errors.password ? 'red.300' : 'gray.300'}
              focusBorderColor={primaryColor}
            />
            {errors.password && <Box color="red.500">{errors.password}</Box>}
          </FormControl>
          {errors.login && (
            <Box color="red.500" textAlign="center">
              {errors.login}
            </Box>
          )}
          <Button
            colorScheme={primaryColor}
            bg={primaryColor}
            width="full"
            color="white"
            onClick={handleSubmit}
          >
            Login
          </Button>
        </Stack>
      </Box>
    </Box>
  );
};

export default Auth;
