import React, { useEffect, useRef, useState, useMemo } from 'react';
import { Box, useTheme, useToast } from '@chakra-ui/react';
import isEqual from 'lodash.isequal';
import { useParams } from 'react-router-dom';
import {
  useAccounts,
  usePayees,
  useCategories,
  useTransactions
} from '../context';
import {
  AccountHeader,
  AddTransactionModal,
  TransactionsTable,
  Loading,
  Error
} from '../components';
import { formatCurrency } from '../utils';

const Account = () => {
  const { accountId } = useParams();
  const toast = useToast();
  const theme = useTheme();
  const primaryColor = theme.colors.primary;
  const {
    transactions,
    loading,
    error,
    getTransactions,
    importTransactions,
    updateTransaction,
    deleteTransaction
  } = useTransactions();
  const { accounts } = useAccounts();
  const { payees } = usePayees();
  const { categories } = useCategories();
  const [isModalOpen, setModalOpen] = useState(false);
  const [localTransactions, setLocalTransactions] = useState(transactions);
  const [originalTransactions, setOriginalTransactions] =
    useState(transactions);
  const [validationErrors, setValidationErrors] = useState({});
  const [, setSelectedFile] = useState(null);
  const fileInputRef = useRef(null);

  useEffect(() => {
    const fetchTransactions = async () => {
      if (accountId) {
        setLocalTransactions([]);
        setOriginalTransactions([]);
        await getTransactions(accountId);
      }
    };

    fetchTransactions();
  }, [accountId, getTransactions]);

  useEffect(() => {
    setLocalTransactions(transactions);
    setOriginalTransactions(JSON.parse(JSON.stringify(transactions)));
  }, [transactions]);

  const totalClearedCredit = useMemo(() => {
    return localTransactions
      .filter(transaction => transaction.clearedAt)
      .reduce((acc, transaction) => acc + (transaction.credit || 0), 0);
  }, [localTransactions]);

  const totalClearedDebit = useMemo(() => {
    return localTransactions
      .filter(transaction => transaction.clearedAt)
      .reduce((acc, transaction) => acc + (transaction.debit || 0), 0);
  }, [localTransactions]);

  if (loading) return <Loading />;
  if (error) return <Error message={error.message} />;

  const handleDelete = async id => {
    await deleteTransaction(id);
    toast({
      title: 'Transaction deleted.',
      description: `Transaction has been successfully deleted.`,
      status: 'success',
      duration: 1500,
      isClosable: true
    });
  };

  const handleInputChange = (index, field, value) => {
    const updatedTransactions = [...localTransactions];
    updatedTransactions[index][field] = value;
    setLocalTransactions(updatedTransactions);

    if (field === 'name' && value) {
      setValidationErrors(prev => ({ ...prev, [index]: false }));
    }
  };

  const handleCheckboxChange = index => {
    const updatedTransactions = [...localTransactions];
    const currentDateTime = new Date().toISOString();
    updatedTransactions[index].clearedAt = updatedTransactions[index].clearedAt
      ? null
      : currentDateTime;
    setLocalTransactions(updatedTransactions);
  };

  const handleSaveChanges = async index => {
    const transaction = localTransactions[index];
    const originalTransaction = originalTransactions[index];

    const errors = {};
    if (!transaction.name) {
      errors.name = true;
    }

    if (Object.keys(errors).length > 0) {
      setValidationErrors(prev => ({ ...prev, [index]: errors }));
      return;
    }

    setValidationErrors(prev => ({ ...prev, [index]: {} }));

    if (transaction.payeeId === '') {
      transaction.payeeId = null;
    }

    const hasChanges = !isEqual(transaction, originalTransaction);
    if (hasChanges) {
      await updateTransaction(transaction.id, {
        ...transaction,
        clearedAt: transaction.clearedAt
      });
    }
  };

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
      const response = await importTransactions(file);
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
        totalClearedCredit={totalClearedCredit}
        totalClearedDebit={totalClearedDebit}
        formatCurrency={formatCurrency}
        primaryColor={primaryColor}
        setModalOpen={setModalOpen}
        openFileDialog={openFileDialog}
        handleFileChange={handleFileChange}
        fileInputRef={fileInputRef}
      />
      <Box my="4" />
      <TransactionsTable
        localTransactions={localTransactions}
        handleInputChange={handleInputChange}
        handleSaveChanges={handleSaveChanges}
        validationErrors={validationErrors}
        payees={payees}
        categories={categories}
        handleCheckboxChange={handleCheckboxChange}
        handleDelete={handleDelete}
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
