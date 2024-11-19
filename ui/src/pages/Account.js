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
  const { currentAccount, getAccountById } = useAccounts();
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

  useEffect(() => {
    const fetchAccount = async () => {
      if (accountId) {
        await getAccountById(accountId);
      }
    };
    fetchAccount();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [accountId]);

  const { name, balance } = currentAccount || {};

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
      await getAccountById(accountId);

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
    } finally {
      setSelectedFile(null);
      fileInputRef.current.value = null;
    }
  };

  const handleFileChange = async event => {
    const file = event.target.files[0];
    if (!file) return;

    await handleFileUpload(file);

    setSelectedFile(null);
    fileInputRef.current.value = null;
  };

  const openFileDialog = () => {
    fileInputRef.current.click();
  };

  return (
    <Box>
      <AccountHeader
        accountName={name}
        accountBalance={balance}
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
      <TransactionsTable getAccountById={getAccountById} />
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
