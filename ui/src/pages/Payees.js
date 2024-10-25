import React, { useState, useEffect } from 'react';
import { Box, useToast } from '@chakra-ui/react';
import { usePayees } from '../context';
import { PayeeTable, Loading, Error } from '../components';

const Payees = () => {
  const { payees, deletePayee, updatePayee, createPayee, loading, error } =
    usePayees();
  const [localPayees, setLocalPayees] = useState([]);
  const [newPayee, setNewPayee] = useState('');
  const [errors, setErrors] = useState({});
  const toast = useToast();

  useEffect(() => {
    setLocalPayees(payees);
  }, [payees]);

  const handleAddPayee = async () => {
    if (newPayee.trim()) {
      try {
        const newPayeeObj = { name: newPayee.trim() };
        await createPayee(newPayeeObj);
        setNewPayee('');
        toast({
          title: 'Payee added.',
          description: `Payee ${newPayee} has been successfully added.`,
          status: 'success',
          duration: 1500,
          isClosable: true
        });
      } catch (err) {
        toast({
          title: 'Error adding payee.',
          description: err.message,
          status: 'error',
          duration: 1500,
          isClosable: true
        });
      }
    } else {
      toast({
        title: 'Error adding payee.',
        description: 'Payee name cannot be empty.',
        status: 'error',
        duration: 1500,
        isClosable: true
      });
    }
  };

  const handleDelete = async id => {
    await deletePayee(id);
    toast({
      title: 'Payee deleted.',
      description: `Payee has been successfully deleted.`,
      status: 'info',
      duration: 1500,
      isClosable: true
    });
  };

  const handleFieldChange = (id, field, value) => {
    const updatedPayees = localPayees.map(payee =>
      payee.id === id ? { ...payee, [field]: value } : payee
    );
    setLocalPayees(updatedPayees);

    setErrors(prevErrors => ({
      ...prevErrors,
      [id]: {
        ...prevErrors[id],
        [field]: !value
      }
    }));
  };

  const handleSaveChanges = async (id, updatedPayee) => {
    if (!updatedPayee.name.trim()) return;

    try {
      await updatePayee(id, updatedPayee);
    } catch (err) {
      toast({
        title: 'Error updating payee.',
        description: err.message,
        status: 'error',
        duration: 1500,
        isClosable: true
      });
    }
  };

  if (loading) return <Loading />;
  if (error) return <Error message={error.message} />;

  return (
    <Box overflowX="auto" borderRadius="md">
      <PayeeTable
        payees={localPayees}
        newPayee={newPayee}
        setNewPayee={setNewPayee}
        handleAddPayee={handleAddPayee}
        handleDelete={handleDelete}
        handleFieldChange={handleFieldChange}
        handleSaveChanges={handleSaveChanges}
        errors={errors}
      />
    </Box>
  );
};

export default Payees;
