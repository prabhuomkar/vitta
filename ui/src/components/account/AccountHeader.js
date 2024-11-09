import React, { useState, useRef } from 'react';
import {
  Card,
  CardBody,
  Flex,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  Input,
  Button
} from '@chakra-ui/react';
import { AddIcon, ArrowDownIcon } from '@chakra-ui/icons';

const AccountHeader = ({
  accountName,
  totalClearedCredit,
  totalClearedDebit,
  formatCurrency,
  primaryColor,
  setModalOpen,
  openFileDialog,
  handleFileChange,
  fileInputRef,
  handleSearch,
  searchQuery
}) => {
  const [inputValue, setInputValue] = useState(searchQuery);
  const debounceTimeoutRef = useRef(null);

  const onInputChange = e => {
    const query = e.target.value;
    setInputValue(query);

    if (debounceTimeoutRef.current) {
      clearTimeout(debounceTimeoutRef.current);
    }

    debounceTimeoutRef.current = setTimeout(() => {
      handleSearch(query);
    }, 800);
  };

  return (
    <Card boxShadow="none" width="auto">
      <CardBody>
        <Flex direction={['column', 'row']} alignItems="center" gap={4}>
          <Stat flex={[null]} textAlign={['center', 'left']}>
            <StatLabel>{accountName}</StatLabel>
            <StatNumber>
              {formatCurrency(totalClearedCredit - totalClearedDebit)}
            </StatNumber>
            <StatHelpText>CASH/CHECK IN</StatHelpText>
          </Stat>
          <Input
            placeholder="Search"
            flex={[null, '2']}
            mx={[0, 4]}
            width={['100%', 'auto']}
            maxLength={255}
            value={inputValue}
            onChange={onInputChange}
          />
          <Button
            leftIcon={<AddIcon boxSize={3} />}
            width={['100%', 'auto']}
            onClick={() => setModalOpen(true)}
          >
            Add
          </Button>
          <Button
            leftIcon={<ArrowDownIcon />}
            width={['100%', 'auto']}
            colorScheme={primaryColor}
            bgColor={primaryColor}
            onClick={openFileDialog}
          >
            Import
          </Button>
          <Input
            type="file"
            ref={fileInputRef}
            onChange={handleFileChange}
            accept=".csv, .xls, .xlsx"
            width={['100%', 'auto']}
            display="none"
            maxLength={255}
          />
        </Flex>
      </CardBody>
    </Card>
  );
};

export default AccountHeader;
