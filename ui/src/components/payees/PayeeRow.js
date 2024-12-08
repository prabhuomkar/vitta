import React, { useState } from 'react';
import { Tr, Td, Input, Select, Tooltip, IconButton } from '@chakra-ui/react';
import { DeleteIcon } from '@chakra-ui/icons';
import { FaListUl } from 'react-icons/fa';
import SetRulesModal from './modal/SetRulesModal';

const PayeeRow = ({
  payee,
  categories,
  handleDelete,
  handleFieldChange,
  handleSaveChanges,
  isFieldInvalid
}) => {
  const [rules, setRules] = useState(payee.rules || {});
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleCategoryChange = e => {
    const autoCategoryId = e.target.value || null;
    handleFieldChange(payee.id, 'autoCategoryId', autoCategoryId);
  };

  const handleCategoryBlur = () => {
    handleSaveChanges(payee.id, {
      ...payee,
      autoCategoryId: payee.autoCategoryId
    });
  };

  const handleSaveRules = async updatedRules => {
    try {
      setRules(updatedRules);
      handleFieldChange(payee.id, 'rules', updatedRules);

      const updatedPayee = {
        ...payee,
        rules: updatedRules
      };
      await handleSaveChanges(payee.id, updatedPayee);
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Error updating rules:', error.message);
    }
  };

  const openModal = () => {
    setIsModalOpen(true);
  };

  return (
    <>
      <Tr>
        <Td padding="0.6rem">
          <Input
            value={payee.name || ''}
            onChange={e => handleFieldChange(payee.id, 'name', e.target.value)}
            onBlur={() =>
              handleSaveChanges(payee.id, { ...payee, name: payee.name })
            }
            placeholder="Payee Name"
            isInvalid={isFieldInvalid}
            errorBorderColor="red.500"
            size="sm"
            maxLength={255}
          />
        </Td>
        <Td padding="0.6rem">
          <Select
            value={payee.autoCategoryId || ''}
            onChange={handleCategoryChange}
            onBlur={handleCategoryBlur}
            placeholder="Select Category"
            size="sm"
          >
            {categories.map(category => (
              <option key={category.id} value={category.id}>
                {category.name}
              </option>
            ))}
          </Select>
        </Td>
        <Td padding="0.6rem">
          <Tooltip label="Set Rules" hasArrow openDelay={500}>
            <IconButton
              icon={<FaListUl />}
              variant="outline"
              size="sm"
              onClick={openModal}
            />
          </Tooltip>
        </Td>
        <Td padding="0.6rem">
          <Tooltip label="Delete Payee" hasArrow openDelay={500}>
            <IconButton
              aria-label="Delete payee"
              icon={<DeleteIcon />}
              variant="outline"
              onClick={() => handleDelete(payee.id)}
              size="sm"
            />
          </Tooltip>
        </Td>
      </Tr>
      <SetRulesModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSave={handleSaveRules}
        initialValues={rules}
      />
    </>
  );
};

export default PayeeRow;
