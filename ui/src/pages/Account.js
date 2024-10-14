import React from 'react';
import {
  Box,
  Card,
  CardBody,
  Flex,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  Input,
  Button,
  useTheme
} from '@chakra-ui/react';
import Transcations from '../components/Transcations';

const Account = () => {
  const theme = useTheme();
  const primaryColor = theme.colors.primary;

  return (
    <>
      <Card boxShadow="none" width="auto">
        <CardBody>
          <Flex direction={['column', 'row']} alignItems="center" gap={4}>
            <Stat flex={[null]} textAlign={['center', 'left']}>
              <StatLabel>XYZ Bank</StatLabel>
              <StatNumber>₹100,000</StatNumber>
              <StatHelpText>CASH/CHECK IN</StatHelpText>
            </Stat>

            <Input
              placeholder="Search"
              flex={[null, '2']}
              mx={[0, 4]}
              width={['100%', 'auto']}
            />

            <Button
              flex={[null]}
              width={['100%', 'auto']}
              colorScheme={primaryColor}
              bgColor={primaryColor}
            >
              Import
            </Button>
          </Flex>
        </CardBody>
      </Card>
      <Box my="4" />
      <Transcations />
    </>
  );
};

export default Account;
