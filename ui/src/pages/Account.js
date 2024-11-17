import React, { useEffect, useRef, useState } from 'react';
import { Box, useTheme, useToast } from '@chakra-ui/react';
import { useParams } from 'react-router-dom';
import { useAccounts, useTransactions } from '../context';
import {
  AccountHeader,
  AddTransactionModal,
  TransactionsTable,
  Pagination,
  Error
} from '../components';
import { formatCurrency } from '../utils';

const Account = () => {
  const { accountId } = useParams();
  const toast = useToast();
  const theme = useTheme();
  const primaryColor = theme.colors.primary;
  const {
    total,
    totalPages,
    loading,
    error,
    getTransactions,
    importTransactions,
    searchQuery,
    updateSearchQuery,
    goToNextPage,
    goToPreviousPage,
    page,
    hasNextPage
  } = useTransactions();
  const { accounts } = useAccounts();
  const [isModalOpen, setModalOpen] = useState(false);
  const [, setSelectedFile] = useState(null);
  const fileInputRef = useRef(null);

  useEffect(() => {
    const fetchTransactions = async () => {
      if (accountId) {
        await getTransactions(accountId);
      }
    };

    fetchTransactions();
  }, [accountId, getTransactions]);

  const handleSearch = query => {
    updateSearchQuery(query);
  };

  if (error) return <Error message={error.message} />;

  const handleFileUpload = async file => {
    if (!file) {
      toast({
        title: 'No file selected',
        description: 'Please select a file to upload.',
        status: 'warning',
        duration: 1500,
        isClosable: true
      });
      return;
    }

    try {
      const response = await importTransactions(accountId, file);
      await getTransactions(accountId);

      if (response?.success) {
        toast({
          title: 'Transactions imported successfully',
          description: `The file has been uploaded, and ${response?.transactions?.imported} transactions are imported.`,
          status: 'success',
          duration: 1500,
          isClosable: true
        });
      }
    } catch (err) {
      toast({
        title: 'Upload failed',
        description: 'There was an error uploading the file. Please try again.',
        status: 'error',
        duration: 1500,
        isClosable: true
      });
    }
  };

  const handleFileChange = event => {
    setSelectedFile(event.target.files[0]);
    handleFileUpload(event.target.files[0]);
  };

  const openFileDialog = () => {
    fileInputRef.current.click();
  };

  const accountName = accounts.find(account => account.id === accountId)?.name;

  return (
    <Box>
      <AccountHeader
        accountName={accountName}
        formatCurrency={formatCurrency}
        primaryColor={primaryColor}
        setModalOpen={setModalOpen}
        openFileDialog={openFileDialog}
        handleFileChange={handleFileChange}
        fileInputRef={fileInputRef}
        handleSearch={handleSearch}
        searchQuery={searchQuery}
        total={total}
      />
      <Box my="4" />
      <TransactionsTable />
      <Pagination
        page={page}
        totalPages={totalPages}
        goToNextPage={goToNextPage}
        goToPreviousPage={goToPreviousPage}
        hasNextPage={hasNextPage}
        loading={loading}
      />
      <AddTransactionModal
        isOpen={isModalOpen}
        onClose={() => setModalOpen(false)}
        getTransactions={() => getTransactions(accountId)}
      />
    </Box>
  );
};

export default Account;
